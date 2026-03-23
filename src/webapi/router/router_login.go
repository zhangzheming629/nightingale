package router

import (
	"fmt"
	"net/http"
	"strings"
        "errors"
        "encoding/json"
        "io"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"

	"github.com/didi/nightingale/v5/src/models"
	"github.com/didi/nightingale/v5/src/webapi/config"
)

type loginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func CallAuth(token string) (string, error){
  address := config.C.AuthServer.Address
  requestUrl := address+"/auth2/api/v2/user/getLoginUser"
  fmt.Println("CallAuth requestUrl:", requestUrl) 
  req,_ := http.NewRequest("GET", requestUrl, nil) 
  fmt.Println("token:", token)
  //token = "Bearer " + token 
  req.Header.Add("Authorization", token)
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil{
    fmt.Println(err)
    return "",err
  }
  resBody, _ := io.ReadAll(resp.Body)
  requestBody := string(resBody)
  if requestBody == "400 Bad Request" {
     return requestBody, errors.New("请求auth错误")
  }
  defer resp.Body.Close()
  return requestBody, nil
}

type AuthReturnBody struct {
  Status int `json:"status"`
  Msg    string `json:"msg"`
  Success bool `json:"success"`
  Data   Data `json:"data"`
}

type Data struct{
  // UserId string `json:"userId"`
  // TenantName string `json: "tenantName"`
  // TenantId string `json: "tenantId"`
  UserLoginName string `json: "userLoginName"`
  PhoneNumber string `json: "phoneNumber"`
  EmailAddress string `json: "emailAddress"` 
}


type AuthLoginReturnBody struct {
  Status int `json:"status"`
  Msg    string `json:"msg"`
  Success bool `json:"success"`
  Data   Data `json:"data"`
}


func authLoginPost(c *gin.Context) {
     fmt.Println("authLoginPost")
     token := c.GetHeader("Authorization")
     fmt.Println(token)
     authReturn, err := CallAuth(token)     
     if err != nil {
        fmt.Println(err)
        ginx.NewRender(c).Message(err)
        return
     }
     fmt.Println(authReturn)
     var result AuthReturnBody 
     err = json.Unmarshal([]byte(authReturn), &result)
     if err != nil {
       fmt.Println(err)       
       returnBody := AuthLoginReturnBody{
         Status: result.Status,
         Msg: result.Msg,
         Success: result.Success,
         Data: result.Data,
       }
       ginx.NewRender(c).Data(returnBody, nil)
       return
     } 
    
     fmt.Println(result)
     fmt.Println(result.Data)
     fmt.Println(result.Data.UserLoginName)
     fmt.Println(result.Data.PhoneNumber)
     fmt.Println(result.Data.EmailAddress)
     
     // user, err := UserGetByUsername(u.Username)
     returnBody := AuthLoginReturnBody{
       Status: result.Status,
       Msg: result.Msg,
       Success: result.Success,
       Data: result.Data,
     }
     ginx.NewRender(c).Data(returnBody, nil)
}

func loginPost(c *gin.Context) {
	var f loginForm
	ginx.BindJSON(c, &f)

	user, err := models.PassLogin(f.Username, f.Password)
	if err != nil {
		// pass validate fail, try ldap
		if config.C.LDAP.Enable {
			user, err = models.LdapLogin(f.Username, f.Password)
			if err != nil {
				ginx.NewRender(c).Message(err)
				return
			}
		} else {
			ginx.NewRender(c).Message(err)
			return
		}
	}

	if user == nil {
		// Theoretically impossible
		ginx.NewRender(c).Message("Username or password invalid")
		return
	}

	userIdentity := fmt.Sprintf("%d-%s", user.Id, user.Username)

	ts, err := createTokens(config.C.JWTAuth.SigningKey, userIdentity)
	ginx.Dangerous(err)
	ginx.Dangerous(createAuth(c.Request.Context(), userIdentity, ts))

	ginx.NewRender(c).Data(gin.H{
		"user":          user,
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}, nil)
}

func logoutPost(c *gin.Context) {
	metadata, err := extractTokenMetadata(c.Request)
	if err != nil {
		ginx.NewRender(c, http.StatusBadRequest).Message("failed to parse jwt token")
		return
	}

	delErr := deleteTokens(c.Request.Context(), metadata)
	if delErr != nil {
		ginx.NewRender(c).Message(InternalServerError)
		return
	}

	ginx.NewRender(c).Message("")
}

type refreshForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func refreshPost(c *gin.Context) {
	var f refreshForm
	ginx.BindJSON(c, &f)

	// verify the token
	token, err := jwt.Parse(f.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
		}
		return []byte(config.C.JWTAuth.SigningKey), nil
	})

	// if there is an error, the token must have expired
	if err != nil {
		// redirect to login page
		ginx.NewRender(c, http.StatusUnauthorized).Message("refresh token expired")
		return
	}

	// Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			// Theoretically impossible
			ginx.NewRender(c, http.StatusUnauthorized).Message("failed to parse refresh_uuid from jwt")
			return
		}

		userIdentity, ok := claims["user_identity"].(string)
		if !ok {
			// Theoretically impossible
			ginx.NewRender(c, http.StatusUnauthorized).Message("failed to parse user_identity from jwt")
			return
		}

		// Delete the previous Refresh Token
		err = deleteAuth(c.Request.Context(), refreshUuid)
		if err != nil {
			ginx.NewRender(c, http.StatusUnauthorized).Message(InternalServerError)
			return
		}

		// Delete previous Access Token
		deleteAuth(c.Request.Context(), strings.Split(refreshUuid, "++")[0])

		// Create new pairs of refresh and access tokens
		ts, err := createTokens(config.C.JWTAuth.SigningKey, userIdentity)
		ginx.Dangerous(err)
		ginx.Dangerous(createAuth(c.Request.Context(), userIdentity, ts))

		ginx.NewRender(c).Data(gin.H{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}, nil)
	} else {
		// redirect to login page
		ginx.NewRender(c, http.StatusUnauthorized).Message("refresh token expired")
	}
}

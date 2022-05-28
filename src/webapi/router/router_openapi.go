package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"strconv"
	// "os/exec"
	// "strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// "encoding/json"
	// "os"
	"github.com/didi/nightingale/v5/src/service"
)

type ClusterInfo struct {
	Scdcs Scdcs `json:"scdcs"`
	Coordinator Coordinator `json:"coordinator"`
	Executor Executor `json:"executor"`
	Gtm Gtm `json:"gtm"`
	Storage Storage `json:"storage"`
	BackupHost BackupHost `json:"backup_host"`
	Setup Setup `json:"setup"`
	Type string `json:"type"`
	Address string `json:"address"`
	MirrorAdress string `json:"mirror_adress"`
	ServicePort int `json:"service_port"`
	InnerPort int `json:"inner_port"`
	MirrorPort int `json:"mirror_port"`
	Count string `json:"count"`
	MirrorHost string `json:"mirror_host"`
	MirrorDataDir string `json:"mirror_data_dir"`
}

type Setup struct {
	Hostname string `json:"hostname"`
	Address string `json:"address"`
}
type Scdcs struct {
	BaseServicePort int `json:"base_service_port"`
	BaseInnerPort int `json:"base_inner_port"`
	BaseDataDir string `json:"base_data_dir"`
	Setup Setup `json:"setup"`
}

type Coordinator struct {
	BasePort int `json:"base_port"`
	BaseDataDir string `json:"base_data_dir"`
	Setup Setup `json:"setup"`
}

type Executor struct {
	HaveMirror string `json:"have_mirror"`
	MirrorMode string `json:"mirror_mode"`
	BasePort int `json:"base_port"`
	BaseMirrorPort int `json:"base_mirror_port"`
	BaseDataDir string `json:"base_data_dir"`
	MirrorDataDir string `json:"mirror_data_dir"`
	ExecutorCountEachHost int `json:"executor_count_each_host"`
	Setup Setup `json:"setup"`
}

type Gtm struct {
	Port int `json:"port"`
	DataDir string `json:"data_dir"`
	Hostname string `json:"hostname"`
}

type Monitor struct {
	Port int `json:"port"`
	DataDir string `json:"data_dir"`
	Hostname string `json:"hostname"`
	Role string `json:"role"`
}

type Storage struct {
	Type string `json:"type"`
}

type BackupHost struct {
}

type UninstallBody struct {
	Hostname string `json:"hostname"`
	ForceClearData string `json:"force_clear_data"`
	ClusterId string `json:"cluster_id"`
}


func installCluster(c *gin.Context){
	fmt.Println(c)
	fmt.Println("installCluster")
	var status int = 1
	var msg string = "install success"
	var setup Setup
	var err error
	c.BindJSON(&setup)
	fmt.Println("==data===>>>", setup)
	fmt.Println("==data===>>>", setup.Hostname)
	
	// var output []byte
	// var  whoami []byte
	// var cmd *exec.Cmd 
	// cmd = exec.Command("sh", "-c", "/home/zhangzheming/Documents/goCode/src/test.sh")
	// if output, err = cmd.Output(); err != nil {
	// 	fmt.Println(err)
	// }
	// result := string(output)
	// fmt.Println(result)

	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type PostgresConnect struct {
	ConnectAddress string 
	ConnectUser string 
	ConnectPassword string 
	ConnectDatabase string 
	ConnectSslmode bool
}

type CreateDatabaseBody struct {
	Host string `json:"host"`
	Port  int `json:"port"`
	ConnectUser string `json:"connect_user"`
	ConnectPassword string   `json:"connect_password"`
	ConnectDatabase string  `json:"connect_database"`
	ConnectSslmode  bool		 `json:"connect_sslmode`
}

func Pginit(connect PostgresConnect) service.Service {
	svc := service.Service{
		ConnectAddress: connect.ConnectAddress,
		ConnectUser: connect.ConnectUser, 
		ConnectPassword: connect.ConnectPassword,
		ConnectDatabase: connect.ConnectDatabase,
		ConnectSslmode: connect.ConnectSslmode,
		MaxIdle: 1,
		MaxOpen: 1,
	}
	return svc
}
func initCluster(c *gin.Context){
	var status int = 1
	var msg string = "initizlization successed"
	var cluster_id string = "seaboxmpp-e4230eba4894"
	var clusterInfo ClusterInfo
	var err error
	c.BindJSON(&clusterInfo)
	fmt.Println("==data===>>>", clusterInfo)


	b, err:=yaml.Marshal(clusterInfo)
	if err != nil {
		fmt.Println("error:", err)
	}

	err =ioutil.WriteFile("cluster/test.yaml", b ,0666)
	if err != nil {
		fmt.Println("error:", err)
	}
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
		"cluster_id":  cluster_id,
	})
}

func createDatabase(c *gin.Context){
	var status int = 1
	var msg string = "initizlization successed"
	var cluster_id string = "seaboxmpp-e4230eba4894"
	var body CreateDatabaseBody
	var err error
	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)

	var query  = `create database zzm`
	var address string = body.Host + ":" + strconv.Itoa(body.Port)
	connect := PostgresConnect {
		ConnectAddress: address,
		ConnectUser: body.ConnectUser,
		ConnectPassword: body.ConnectPassword,
		ConnectDatabase: body.ConnectDatabase,
		ConnectSslmode: body.ConnectSslmode,

	}
	pg := Pginit(connect)
	err = pg.ExecQuery(query)

	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
		"cluster_id":  cluster_id,
	})
}

func uninstallCluster(c *gin.Context){
	var status int = 1
	var msg string = "uninstall successed"
	var uninstallBody UninstallBody
	var err error
	c.BindJSON(&uninstallBody)
	fmt.Println("==data===>>>", uninstallBody)
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type LogBody struct {
	LogDir string  `json:"log_dir"`
	ClusterId string  `json:"cluster_id"`
}

func modifyLogDir(c *gin.Context){
	var status int = 1
	var msg string = "modify successfully"
	var logBody LogBody
	var err error
	c.BindJSON(&logBody)
	fmt.Println("==data===>>>", logBody)
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type DumpBody struct {
	File string `json:"file"`
	Format string `json:"format"`
	Jobs string `json:"jobs"`
	LockWaitTimeOut string `json:"lock-wait-timeOut"`
	DataOnly string `json:"data_only"`
	Blobs string `json:"blobs"`
	Clean string `json:"clean"`
	Schema string `json:"schema"`
	SchemaOnly string `json:"schema_only"`
	Table string `json:"table"`
	Dbname string `json:"dbname"`
	Host string `json:"host"`
	Port int `json:"port"`
	ClusterId string `json:"cluster_id"`
}
func dumpDatabase(c *gin.Context){
	var status int = 1
	var msg string = "database dump succeeded"
	var dumpBody DumpBody
	var err error
	c.BindJSON(&dumpBody)
	fmt.Println("==data===>>>", dumpBody)
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type ModifyConfigBody struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Remove string `json:"remove"`
	ClusterId string `json:"cluster_id"`
}
func modifyConfig(c *gin.Context){
	var status int = 1
	var msg string = "modiied successfully"
	var modifyConfigBody ModifyConfigBody
	var err error
	c.BindJSON(&modifyConfigBody)
	fmt.Println("==data===>>>", modifyConfigBody)
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type ShowConfigBody struct {
	Name string `json:"name"`
	ClusterId string `json:"cluster_id"`
}
func showConfig(c *gin.Context){
	var status int = 1
	var msg string = "get successfully"
	var showConfigBody ShowConfigBody
	var err error
	var name string = "debug_pretty_print"
    var value string = "on"

	c.BindJSON(&showConfigBody)
	fmt.Println("==data===>>>", showConfigBody)
	ginx.Dangerous(err)
	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
		"name": name,
		"value": value,
	})
}

type DescribeClusterDetailedStatusBody struct {
	Module string `json:"module"`
	Host string `json:"host"`
	ClusterId string `json:"cluster_id"`
}

func describeClusterDetailedStatus(c *gin.Context){
	var status int = 1
	var msg string = "get successfully"
	var body DescribeClusterDetailedStatusBody
	var err error
	

	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)
	ginx.Dangerous(err)


	scdcsSetup := Setup{
		Hostname: "node1",
		Address: "127.0.0.1",
	}
	scdcs := Scdcs{
		BaseServicePort: 38000,
		BaseInnerPort: 38000,
		BaseDataDir: "/data",
		Setup: scdcsSetup,
	}

	coordinatorSetup := Setup{
		Hostname: "node2",
		Address: "127.0.0.1",
	}
	coordinator := Coordinator{
		BasePort: 38000,
		BaseDataDir: "/data",
		Setup: coordinatorSetup,
	}

	executorSetup := Setup{
		Hostname: "node3",
		Address: "127.0.0.1",
	}
	executor := Executor{
		HaveMirror: "true",
		MirrorMode: "auto",
		BasePort: 38000,
		BaseMirrorPort: 38001,
		BaseDataDir: "/data",
		ExecutorCountEachHost: 1,
		Setup: executorSetup,
	}

	gtm := Gtm{
		Hostname: "node1",
		DataDir: "/data",
		Port: 38125,
	}

	monitor := Monitor{
		Hostname: "node1",
		DataDir: "/data",
		Port: 38125,
		Role: "primary",
	}

	if body.Module == "coordinator" {
		c.JSON(200, gin.H{
			"status": status,
			"msg":  msg,

			"coordinator": coordinator,
		})
	} else if body.Module == "scdcs" {
		c.JSON(200, gin.H{
			"status": status,
			"msg":  msg,
			"scdcs": scdcs,
		})
	}else if body.Module == "executor" {
		c.JSON(200, gin.H{
			"status": status,
			"msg":  msg,
			"executor": executor,
		})
	}else if body.Module == "gtm" {
		c.JSON(200, gin.H{
			"status": status,
			"msg":  msg,
			"gtm": gtm,
		})
	}else if body.Module == "monitor" {
		c.JSON(200, gin.H{
			"status": status,
			"msg":  msg,
			"monitor": monitor,
		})
	}
}

type ExpandExecutorBody struct {
	ClusterId string `json:"cluster_id"`
	ExpandExecutor []ExpandExecutor `json:"executor"`
}

type ExpandExecutor struct {
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	DataDir string `json:"data_dir"`
	Role string `json:"role"`
}

func expandExecutor(c *gin.Context){
	var status int = 1
	var msg string = "expand successfully"
	var body ExpandExecutorBody
	var err error

	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)
	ginx.Dangerous(err)

	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}


type ShrinkExecutorBody struct {
	ClusterId string `json:"cluster_id"`
	Count int `json:"count"`
}

func shrinkExecutor(c *gin.Context){
	var status int = 1
	var msg string = "shrinked successfully"
	var body ShrinkExecutorBody
	var err error

	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)
	ginx.Dangerous(err)

	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type StartClusterBody struct {
	ClusterId string `json:"cluster_id"`
}

func startCluster(c *gin.Context){
	var status int = 1
	var msg string = "started successfully"
	var body StartClusterBody
	var err error

	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)
	ginx.Dangerous(err)

	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}

type StopClusterBody struct {
	ClusterId string `json:"cluster_id"`
}

func stopCluster(c *gin.Context){
	var status int = 1
	var msg string = "stopped successfully"
	var body StartClusterBody
	var err error

	c.BindJSON(&body)
	fmt.Println("==data===>>>", body)
	ginx.Dangerous(err)

	c.JSON(200, gin.H{
		"status": status,
		"msg":  msg,
	})
}
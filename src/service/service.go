package service

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"log"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"encoding/base64"
	_ "github.com/lib/pq"
)

type Service struct {
	MaxIdle       int
	MaxOpen       int
	DB            *sql.DB
	ConnectAddress       string
	ConnectUser string 
	ConnectPassword string 
	ConnectDatabase string 
	ConnectSslmode bool
}

var socketRegexp = regexp.MustCompile(`/\.s\.PGSQL\.\d+$`)

func (p *Service) getDb() (db *sql.DB) {
	const localhost = "host=localhost sslmode=disable"
	var err error 
	p.ConnectAddress, err = p.getDbAddress()
	if err != nil {
		log.Println("WARN! getDbAddress Error: ", err)
	}
	if p.ConnectAddress == "" {
		p.ConnectAddress = localhost
	}
	connConfig, err := pgx.ParseConfig(p.ConnectAddress)
	if err != nil {
		log.Println("WARN! ParseConfig Error: ", err)
	}

	// Remove the socket name from the path
	connConfig.Host = socketRegexp.ReplaceAllLiteralString(connConfig.Host, "")
	connectionString := stdlib.RegisterConnConfig(connConfig)
	if p.DB, err = sql.Open("pgx", connectionString); err != nil {
		log.Println("WARN! sql.Open Error: ", err)
	}

	p.DB.SetMaxOpenConns(p.MaxOpen)
	p.DB.SetMaxIdleConns(p.MaxIdle)
	return p.DB
  }

func (p *Service) getDbAddress() (realAddress string, err error) {
	var decodePwdPrefix = "REZKWA=="
	var addressPrefx = "postgres://"
	var dbAddress string 
	var sslmode string 
	var addresses []string 
	if p.ConnectSslmode == false {
		sslmode = "sslmode=disable"
	} else {
		sslmode ="sslmode=verify-ca"
	}
	if strings.Contains(p.ConnectPassword, decodePwdPrefix) {
		p.ConnectPassword = strings.Split(p.ConnectPassword, decodePwdPrefix)[1]
		sDec, err := base64.StdEncoding.DecodeString(p.ConnectPassword)
		if err != nil {
			log.Println("WARN! decode password error: ", err)
		}
		p.ConnectPassword = string(sDec)

	}  
	if !strings.Contains(p.ConnectAddress, addressPrefx){
		addresses = strings.Split(p.ConnectAddress, ",")
		for _,value := range addresses {
			dbAddress = addressPrefx + p.ConnectUser +":"+ p.ConnectPassword +"@"+value+"/"+p.ConnectDatabase+"?"+sslmode
			fmt.Println(dbAddress)
			// postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable
			db, _ := sql.Open("postgres", dbAddress)
			err := db.Ping()
			db.Close()
			if err != nil {
				continue
			} else {
				realAddress = dbAddress
				break
			}
		}
	} else {
		realAddress = p.ConnectAddress
	}
	if realAddress ==""{
		return realAddress, fmt.Errorf("can not connect the db address: %s", addresses)
	}
	return realAddress, err
}

func (p *Service) ExecQuery(query string) error {
	rows, err := p.getDb().Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
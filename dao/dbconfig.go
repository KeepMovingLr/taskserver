package dao

import (
	"database/sql"
	"log"
)

var (
	DB *sql.DB
)

// to reduce the numbers of parameter
type MysqlInitParam struct {
	User         string
	Password     string
	Host         string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
}

func ConnectMysql(param MysqlInitParam) (err error) {
	DB, err = sql.Open("mysql", param.User+":"+param.Password+"@("+param.Host+")/"+param.DBName+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
	// need more test to get the best value
	DB.SetMaxOpenConns(param.MaxOpenConns)
	DB.SetMaxIdleConns(param.MaxIdleConns)
	// each time finish initialize the db, print a log
	log.Println("init database end")
	return nil
}

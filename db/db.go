package db

import (
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Start the connection to database
func Init() *gorm.DB {
	//DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(mysql.Open(CONNECT), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database : error=%v", err)
		return nil
	}

	return db
}

func GetDB() *gorm.DB {
	return db
}

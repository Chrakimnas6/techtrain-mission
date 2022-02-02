package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID    uint   `gorm:"AUTO_INCREMENT"`
	Name  string `json:"name"`
	Token string
}

func main() {
	db := sqlConnect()
	db.AutoMigrate(&User{})
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		db := sqlConnect()
		var users []User
		db.Find(&users)
		defer db.Close()

		c.HTML(200, "index.html", gin.H{
			"users": users,
		})
	})

	router.POST("/user/create", func(c *gin.Context) {
		db := sqlConnect()
		var user User
		c.BindJSON(&user)
		// //generate hash
		// h := sha1.New()
		// s := user.Name + time.Now().String()
		// fmt.Println(s)
		// h.Write([]byte(s))
		// token := h.Sum(nil)
		token := uuid.New()

		user.Token = string(token.String())

		fmt.Println("create user " + user.Name + " with token " + user.Token)
		db.Create(&user)

		//redirect to the home page
		c.Redirect(302, "/")
	})

	router.GET("/user/get", func(c *gin.Context) {
		db := sqlConnect()
		defer db.Close()
		//token := c.GetHeader("x-token")

	})

	router.Run()

}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				fmt.Println("Connection failed")
				panic(err)
			}
			db, err = gorm.Open(DBMS, CONNECT)
		}
	}
	fmt.Println("Connection succeeded")

	return db
}

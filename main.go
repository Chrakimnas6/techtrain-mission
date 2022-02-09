package main

import (
	"training/controllers"
	"training/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := setupServer()
	r.Run()
}

func setupServer() *gin.Engine {
	r := gin.Default()
	// Showing contents of the database
	r.LoadHTMLGlob("templates/*.html")

	db := db.Init()
	controller := controllers.New(db)
	corsConfig := cors.DefaultConfig()

	// CORS setting in order to receive from Swagger properly
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// APIs - user
	r.GET("/", controller.GetUsers)
	r.POST("/user/create", controller.CreateUser)
	r.GET("/user/get", controller.GetUser)
	r.PUT("/user/update", controller.UpdateUser)

	// APIs - gacha
	// APIs - character

	return r
}

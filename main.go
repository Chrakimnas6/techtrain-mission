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

	// Connect to the database
	db := db.Init()
	controller := controllers.Controller{
		Db: db,
	}
	corsConfig := cors.DefaultConfig()

	// CORS setting in order to receive from Swagger properly
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// APIs - user
	r.GET("/", controller.GetAll)
	r.POST("/user/create", controller.CreateUser)
	r.GET("/user/get", controller.GetUser)
	r.PUT("/user/update", controller.UpdateUser)

	// APIs - gacha
	r.POST("/gacha/draw", controller.DrawGacha)
	// APIs - character
	r.GET("/character/list", controller.GetUserCharacters)

	// Internal Use
	// Create new character
	r.POST("/character/create", controller.CreateCharacter)
	// Update possibility
	r.PUT("/gacha/update", controller.UpdateOdds)
	// Create new gacha pool
	r.POST("/gacha/create", controller.CreateGachaPool)

	return r
}

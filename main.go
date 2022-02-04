package main

import (
	"training/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    uint   `gorm:"AUTO_INCREMENT"`
	Name  string `json:"name"`
	Token string
}

func main() {
	r := setupRouter()
	r.Run()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// Showing contents of the database
	r.LoadHTMLGlob("templates/*.html")

	userRepo := controllers.New()
	corsConfig := cors.DefaultConfig()

	// CORS setting in order to receive from Swagger properly
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	// APIs
	r.GET("/", userRepo.GetUsers)
	r.POST("/user/create", userRepo.CreateUser)
	r.GET("/user/get", userRepo.GetUser)
	r.PUT("/user/update", userRepo.UpdateUser)

	return r
}

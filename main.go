package main

import (
	"log"
	"training/controllers"
	"training/db"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// Set up Ethereum network
	client, err := ethclient.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	// Load Keystore
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)

	controller := controllers.Controller{
		Db:       db,
		Client:   client,
		Keystore: ks,
	}

	// Set up CORS in order to receive from Swagger properly
	corsConfig := cors.DefaultConfig()
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
	// APIS - token
	r.POST("/token/transfer", controller.TransferToken)

	// Internal Use
	// Create new character
	r.POST("/character/create", controller.CreateCharacter)
	// Update possibility
	r.PUT("/gacha/update", controller.UpdateOdds)
	// Create new gacha pool
	r.POST("/gacha/create", controller.CreateGachaPool)
	// Create new admin address
	r.POST("/user/admin/create", controller.CreateAdminUser)
	// Get user's ETH and token balance
	r.GET("/user/balance", controller.GetUserBalance)
	// Deploy token
	r.POST("/token/deploy", controller.DeployToken)
	// Mint extra token
	r.POST("/token/mint", controller.MintToken)
	// Check token
	r.GET("/token/check", controller.CheckToken)

	return r
}

package main

import (
	"fmt"
	"log"
	"training/controllers"
	"training/db"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var KS *keystore.KeyStore

func main() {
	SetupEthereumNetwork()
	r := setupServer()
	r.Run()
}
func init() {
	KS = keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
}

//TODO: Set global client and ks
func SetupEthereumNetwork() {
	// Connect to the Ethereum network
	client, err := ethclient.Dial("http://hardhat:8545/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	_ = client

	// Test
	// ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	// password := "secret"

	// account, err := ks.NewAccount(password)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(account.Address.Hex())

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
	// Create new admin address
	r.POST("/admin/create", controller.CreateAdminUser)
	// Get user'balance
	r.GET("/user/balance", controller.GetUserBalance)

	return r
}

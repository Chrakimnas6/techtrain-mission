package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"training/controllers"
	"training/db"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	SetupEthereumNetwork()
	r := setupServer()
	r.Run()
}

func SetupEthereumNetwork() {
	// Connect to the Ethereum network
	client, err := ethclient.Dial("http://hardhat:8545/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ETH balance of the Hardhat's local node: %s\n", balance)
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

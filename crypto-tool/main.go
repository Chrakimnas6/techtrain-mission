package main

import (
	"crypto-tool/controllers"
	"crypto-tool/db"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	client, err := ethclient.Dial("https://eth-ropsten.alchemyapi.io/v2/7dUV0BkLoNzB8UA6bPM9jgrxIzVmSyxG")
	if err != nil {
		log.Fatal(err)
	}
	db := db.Init()

	controller := controllers.Controller{
		Db:     db,
		Client: client,
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/transaction/confirm", controller.ConfirmTransaction)

	r.Run("0.0.0.0:8888") // listen and serve on

}

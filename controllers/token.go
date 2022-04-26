package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"training/accounts"
	token "training/contracts"
	"training/models"
	"training/repos"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AmountRequest struct {
	Amount int `json:"amount"`
}

type AddressRequest struct {
	Address string `json:"address"`
}

// Transfer token to user
func (controller *Controller) TransferToken(c *gin.Context) {
	var user models.User
	userToken := c.GetHeader("x-token")
	err := repos.GetUser(controller.Db, &user, userToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Get the amount
	var amountRequest AmountRequest
	err = c.BindJSON(&amountRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Get user account
	userAccount := accounts.ImportAccount(controller.Keystore, user.Keystore, "password")

	// Get admin's information
	adminUser := models.User{}
	err = repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = token.TransferToken(controller.Client, controller.Keystore, controller.Instance, adminUser.Keystore, userAccount.Address, amountRequest.Amount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	userBalance := token.GetTokenBalance(controller.Instance, userAccount.Address)

	c.JSON(http.StatusOK, gin.H{
		"balance": userBalance,
	})
}

// Deploy Token
// Assume we only have one type of token
func (controller *Controller) DeployToken(c *gin.Context) {
	// Get admin's information
	adminUser := models.User{}
	err := repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Deploy token
	tokenAddress, instance := token.Deploy(controller.Client, controller.Keystore, adminUser.Keystore)

	// Save token information to the database
	tkn := models.Token{
		Address: tokenAddress.Hex(),
		Symbol:  "MTK",
	}
	err = repos.SaveToken(controller.Db, &tkn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	controller.Instance = instance
	token.CheckInformation(instance)

	c.Status(http.StatusOK)

}

// Mint Token
func (controller *Controller) MintToken(c *gin.Context) {
	// Get admin's information
	adminUser := models.User{}
	err := repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Get the amount
	var amountRequest AmountRequest
	err = c.BindJSON(&amountRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// If the instance is nil, get the instance first
	if controller.Instance == nil {
		tkn := models.Token{}
		err := repos.GetToken(controller.Db, &tkn, 1)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// address := common.HexToAddress("0xD1e685605C02f812D4200A16D6844E354ddCDD3C")
		address := common.HexToAddress(tkn.Address)
		instance, err := token.NewToken(address, controller.Client)
		if err != nil {
			log.Fatal(err)
		}
		controller.Instance = instance
	}
	// Mint token
	err = token.MintToken(controller.Client, controller.Keystore, controller.Instance, adminUser.Keystore, amountRequest.Amount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

func (controller *Controller) CheckToken(c *gin.Context) {
	var addressRequest AddressRequest
	err := c.BindJSON(&addressRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Test now
	address := common.HexToAddress(addressRequest.Address)
	instance, err := token.NewToken(address, controller.Client)
	if err != nil {
		log.Fatal(err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	supply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Token's name: %s\n", name)
	fmt.Printf("Token's symbol: %s\n", symbol)
	fmt.Printf("Token's supply: %s\n", supply)

	c.Status(http.StatusOK)

}

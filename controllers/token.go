package controllers

import (
	"errors"
	"fmt"
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
	fmt.Println("Getting user's information...")
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
	fmt.Println("Getting the amount...")
	var amountRequest AmountRequest
	err = c.BindJSON(&amountRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Get user account
	fmt.Println("Getting user's account...")
	userAccount, err := accounts.ImportAccount(controller.Keystore, user.Keystore, "password")
	_ = err
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }

	// Get admin's information
	fmt.Println("Getting admin's information...")
	adminUser := models.User{}
	err = repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Transfer token
	// If the instance is nil, get the instance first
	if controller.Instance == nil {
		tkn := models.Token{}
		err := repos.GetToken(controller.Db, &tkn, 2)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// address := common.HexToAddress("0xD1e685605C02f812D4200A16D6844E354ddCDD3C")
		address := common.HexToAddress(tkn.Address)
		instance, err := token.NewToken(address, controller.Client)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		controller.Instance = instance
	}

	fmt.Println("Transferring token...")
	tx, err := token.TransferToken(controller.Client, controller.Keystore, controller.Instance, adminUser.Keystore, userAccount.Address, amountRequest.Amount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Wait transaction to be mined
	fmt.Println("Waiting for transaction to be mined...")
	err = token.CheckTransaction(tx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Get user balance
	fmt.Println("Getting user's balance...")
	userBalance, err := token.GetTokenBalance(controller.Instance, userAccount.Address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": userBalance,
	})
}

// Deploy Token
// Assume we only have one type of token
func (controller *Controller) DeployToken(c *gin.Context) {
	// Get admin's information
	fmt.Println("Getting admin's information...")
	adminUser := models.User{}
	err := repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Deploy token
	fmt.Println("Deploying token...")
	tx, tokenAddress, instance, err := token.Deploy(controller.Client, controller.Keystore, adminUser.Keystore)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Wait transaction to be mined
	fmt.Println("Waiting for transaction to be mined...")
	err = token.CheckTransaction(tx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Save token information to the database
	tkn := models.Token{
		Address: tokenAddress.Hex(),
		Symbol:  "MTK",
	}
	fmt.Println("Saving token information to the database...")
	err = repos.SaveToken(controller.Db, &tkn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Set instance
	controller.Instance = instance

	c.JSON(http.StatusOK, gin.H{
		"symbol":  tkn.Symbol,
		"address": tkn.Address,
	})

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
		err := repos.GetToken(controller.Db, &tkn, 2)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// address := common.HexToAddress("0xD1e685605C02f812D4200A16D6844E354ddCDD3C")
		address := common.HexToAddress(tkn.Address)
		instance, err := token.NewToken(address, controller.Client)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		controller.Instance = instance
	}

	// Mint token
	tx, err := token.MintToken(controller.Client, controller.Keystore, controller.Instance, adminUser.Keystore, amountRequest.Amount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Wait transaction to be mined
	err = token.CheckTransaction(tx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Minted Amount": amountRequest.Amount,
	})
}

func (controller *Controller) CheckToken(c *gin.Context) {
	requestedAdress := c.Query("address")

	address := common.HexToAddress(requestedAdress)
	instance, err := token.NewToken(address, controller.Client)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	supply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"name":   name,
		"symbol": symbol,
		"supply": supply,
	})

}

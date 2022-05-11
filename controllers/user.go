package controllers

import (
	"errors"
	"fmt"
	"math/big"

	"training/accounts"
	token "training/contracts"
	"training/models"
	"training/repos"

	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Get all users in the database
func (controller *Controller) GetUsers(c *gin.Context) {
	var users []models.User
	err := repos.GetUsers(controller.Db, &users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"users": users,
	})
}

// Create user with "name"
// Use uuid to generate token for the user
func (controller *Controller) CreateUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	userToken := uuid.New()
	user.Token = userToken.String()

	// Create keystore, as for now just use default password
	fmt.Println("Creating keystore...")
	keystoreFileName, account, err := accounts.CreateAccount(controller.Keystore, "password")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	_ = account
	user.Keystore = keystoreFileName

	// Transfer ETH from admin to the user
	// Get admin's information
	fmt.Println("Getting admin's information...")
	adminUser := models.User{}
	err = repos.GetUserByName(controller.Db, &adminUser, "admin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Transfer 0.1 ETH from admin's address to user's address
	fmt.Println("Transferring ETH...")
	tx, err := token.TransferETH(controller.Client, controller.Keystore, adminUser.Keystore, account.Address, 1)
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

	// Create the user
	fmt.Println("Creating user...")
	err = repos.CreateUser(controller.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

// Get user by the token from Header
func (controller *Controller) GetUser(c *gin.Context) {
	var user models.User
	token := c.GetHeader("x-token")
	err := repos.GetUser(controller.Db, &user, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name": user.Name,
	})
}

// Update user by the token from Header with new name
func (controller *Controller) UpdateUser(c *gin.Context) {
	var user models.User
	token := c.GetHeader("x-token")
	err := repos.GetUser(controller.Db, &user, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = repos.UpdateUser(controller.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Status(http.StatusOK)
}

// Get both users and characters from the database
func (controller *Controller) GetAll(c *gin.Context) {
	var users []models.User
	var characters []models.Character
	var characterOdds []models.GachaCharacterOdds
	err := repos.GetUsers(controller.Db, &users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.GetCharacters(controller.Db, &characters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.GetGachaCharacterOddsAll(controller.Db, &characterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"users":                users,
		"characters":           characters,
		"gacha_character_odds": characterOdds,
	})
}

// Get user's balance
func (controller *Controller) GetUserBalance(c *gin.Context) {
	fmt.Println("Getting user's from the database...")
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

	account, err := accounts.ImportAccount(controller.Keystore, user.Keystore, "password")
	_ = err
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }
	var etherBalance, tokenBalance *big.Int
	etherBalance, err = token.GetETHBalance(controller.Client, account.Address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if controller.Instance == nil {
		tkn := models.Token{}
		err := repos.GetToken(controller.Db, &tkn, "MTK")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		address := common.HexToAddress(tkn.Address)
		instance, err := token.NewToken(address, controller.Client)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		controller.Instance = instance
	}
	tokenBalance, err = token.GetTokenBalance(controller.Instance, account.Address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ether balance": etherBalance,
		"token balance": tokenBalance,
	})
}

// Create admin user
func (controller *Controller) CreateAdminUser(c *gin.Context) {
	var user models.User

	// Give the name with admin
	user.Name = "admin"
	userToken := uuid.New()
	user.Token = userToken.String()

	// Check if admin already exist
	err := repos.GetUserByName(controller.Db, &user, "admin")
	if err == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Admin already exists!",
		})
		return
	}

	// Use the existing account as the admin and store it in the path
	account, err := accounts.StoreKey(controller.Keystore, "password")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	user.Keystore = account.URL.Path

	err = repos.CreateUser(controller.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
}

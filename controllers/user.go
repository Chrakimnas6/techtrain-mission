package controllers

import (
	"errors"

	"training/accounts"
	token "training/contracts"
	"training/models"
	"training/repos"

	"net/http"

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
	token := uuid.New()
	user.Token = token.String()

	// Create keystore, at now just use default password
	ks, keystoreFileName, account := accounts.CreateKeystore("password")
	_, _ = ks, account
	user.Keystore = keystoreFileName

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
	c.HTML(200, "index.html", gin.H{
		"users":                users,
		"characters":           characters,
		"gacha_character_odds": characterOdds,
	})
}

// Get user's balance
func (controller *Controller) GetUserBalance(c *gin.Context) {
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

	ks, account := accounts.ImportKeystore(user.Keystore, "password")
	_ = ks
	balance := accounts.GetBalance(account.Address)
	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

// Create admin user
func (controller *Controller) CreateAdminUser(c *gin.Context) {
	var user models.User
	//err := c.BindJSON(&user)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
	// 	return
	// }

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
	// Create keystore, at now just use default password
	ks, keystoreFileName, account := accounts.CreateKeystore("password")
	_ = ks
	user.Keystore = keystoreFileName
	address := account.Address

	err = repos.CreateUser(controller.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Transfer ETH to the admin
	err = token.FaucetTransfer(address)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Deploy the contract
	tokenAddress, instance := token.Deploy(keystoreFileName)
	_ = tokenAddress
	token.CheckInformation(instance, address)

	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

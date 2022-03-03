package controllers

import (
	"errors"

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

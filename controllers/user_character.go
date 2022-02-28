package controllers

import (
	"errors"
	"net/http"

	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get user characters with specfic user's token
func (controller *Controller) GetUserCharacters(c *gin.Context) {
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
	// Get all user_characters of certain user from database
	var userCharactersResponses []models.UserCharacterResponse
	err = repos.GetUserCharacters(controller.Db, &userCharactersResponses, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

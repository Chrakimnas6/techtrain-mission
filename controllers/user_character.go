package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get user characters with specfic user's token
func (controller *Controller) GetUserCharacters(c *gin.Context) {
	var user models.User
	var offset int
	var limit int

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
	// Pagination
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr == "" && limitStr == "" {
		offset = -1
		limit = -1
	} else {
		page, _ := strconv.Atoi(pageStr)
		if page == 0 {
			page = 1
		}
		limit, _ = strconv.Atoi(limitStr)
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}
		offset = (page - 1) * limit
	}
	// Get all user_characters of certain user from database
	var userCharactersResponses []models.UserCharacterResponse
	err = repos.GetUserCharacters(controller.Db, &userCharactersResponses, limit, offset, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

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

type UserCharacterResponse struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
}

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
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		}
		if page == 0 {
			page = 1
		}
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		}
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}
		offset = (page - 1) * limit
	}
	// Get all user_characters of certain user from database
	userCharactersResponses := make([]UserCharacterResponse, 0)
	userCharacters := make([]models.UserCharacter, 0)

	err = repos.GetUserCharacters(controller.Db, &userCharacters, limit, offset, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	for _, userCharacter := range userCharacters {
		userCharactersResponses = append(userCharactersResponses, UserCharacterResponse{
			UserCharacterID: strconv.FormatUint(uint64(userCharacter.ID), 10),
			CharacterID:     strconv.FormatUint(uint64(userCharacter.CharacterID), 10),
			Name:            userCharacter.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

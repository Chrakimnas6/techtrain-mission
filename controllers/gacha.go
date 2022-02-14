package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Draw with times get from request's body
func (controller *Controller) DrawGacha(c *gin.Context) {
	// Get user with the token from Header
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
	// Get gacha times from JSON body
	var gacha models.Gacha
	err = c.BindJSON(&gacha)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Start the gacha, get number of characters
	// Suppose all characters have same odds
	var gachaResultResponse []models.GachaResultResponse
	// Get number of characters
	var size int64
	err = repos.GetSize(controller.Db, &size)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	for i := 0; i < int(gacha.Times); i++ {
		// Randomly pick a character
		characterID := rand.Intn(int(size)) + 1
		userCharacter := models.UserCharacter{UserID: user.ID, CharacterID: uint(characterID)}
		err = repos.CreateUserCharacter(controller.Db, &userCharacter)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		var character models.Character
		err = repos.GetCharacter(controller.Db, &character, uint(characterID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		gachaResultResponse = append(gachaResultResponse, models.GachaResultResponse{CharacterID: strconv.Itoa(characterID), Name: character.Name})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": gachaResultResponse,
	})

}

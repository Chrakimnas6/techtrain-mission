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
	// Suppose the probability of getting SSR, SR and R is 1 : 10 : 89
	// Generate a number from 1 to 100
	var gachaResultResponse []models.GachaResultResponse
	for i := 0; i < int(gacha.Times); i++ {
		var characterID uint
		n := rand.Intn(100) + 1
		if n == 1 {
			err = repos.GetRandCharacter(controller.Db, &characterID, "ssr")
		} else if n <= 11 {
			err = repos.GetRandCharacter(controller.Db, &characterID, "sr")
		} else {
			err = repos.GetRandCharacter(controller.Db, &characterID, "r")
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		// Create user_character
		userCharacter := models.UserCharacter{UserID: user.ID, CharacterID: characterID}
		err = repos.CreateUserCharacter(controller.Db, &userCharacter)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		var character models.Character
		err = repos.GetCharacter(controller.Db, &character, characterID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		gachaResultResponse = append(gachaResultResponse, models.GachaResultResponse{CharacterID: strconv.Itoa(int(characterID)), Name: character.Name})
	}

	c.JSON(http.StatusOK, gin.H{
		"results": gachaResultResponse,
	})

}

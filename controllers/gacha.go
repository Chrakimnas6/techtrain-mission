package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Draw with times get from request's body
// TODO: move parts of the function to helper
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
	var ssrSize, srSize, rSize int64
	// Create map to record how many times a character appears.
	ssrIDToNumber, srIDToNumber, rIDToNumber := make(map[int]int), make(map[int]int), make(map[int]int)

	// Get current number of SSR, SR, R cards
	err = repos.GetSSRSize(controller.Db, &ssrSize)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	err = repos.GetSRSize(controller.Db, &srSize)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	err = repos.GetRSize(controller.Db, &rSize)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	for i := 0; i < int(gacha.Times); i++ {
		n := rand.Intn(100) + 1
		// Get number of different types of cards
		if n == 1 {
			//err = repos.GetRandCharacter(controller.Db, &characterID, "ssr")
			ssrID := rand.Intn(int(ssrSize)) + 1
			ssrIDToNumber[ssrID]++

		} else if n <= 11 {
			//err = repos.GetRandCharacter(controller.Db, &characterID, "sr")
			srID := rand.Intn(int(srSize)) + 1
			srIDToNumber[srID]++
		} else {
			//err = repos.GetRandCharacter(controller.Db, &characterID, "r")
			rID := rand.Intn(int(rSize)) + 1
			rIDToNumber[rID]++
		}
	}

	// TODO: Create user_character
	//var userCharacters []models.UserCharacter
	for ssrID, numbers := range ssrIDToNumber {
		var ssrCharacter models.CharacterSSR
		err = repos.GetSSR(controller.Db, &ssrCharacter, uint(ssrID))
		userCharacters := make([]models.UserCharacter, numbers)
		for i := range userCharacters {
			userCharacters[i] = models.UserCharacter{UserID: user.ID, CharacterID: ssrCharacter.CharacterID}
		}
		err = repos.CreateUserCharacters(controller.Db, &userCharacters)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	for srID, numbers := range srIDToNumber {
		var srCharacter models.CharacterSR
		err = repos.GetSR(controller.Db, &srCharacter, uint(srID))
		userCharacters := make([]models.UserCharacter, numbers)
		for i := range userCharacters {
			userCharacters[i] = models.UserCharacter{UserID: user.ID, CharacterID: srCharacter.CharacterID}
		}
		err = repos.CreateUserCharacters(controller.Db, &userCharacters)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	for rID, numbers := range rIDToNumber {
		var rCharacter models.CharacterR
		err = repos.GetR(controller.Db, &rCharacter, uint(rID))
		userCharacters := make([]models.UserCharacter, numbers)
		for i := range userCharacters {
			userCharacters[i] = models.UserCharacter{UserID: user.ID, CharacterID: rCharacter.CharacterID}
		}
		err = repos.CreateUserCharacters(controller.Db, &userCharacters)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
	// userCharacter := models.UserCharacter{UserID: user.ID, CharacterID: characterID}
	// err = repos.CreateUserCharacter(controller.Db, &userCharacter)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }
	// var character models.Character
	// err = repos.GetCharacter(controller.Db, &character, characterID)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }

	// gachaResultResponse = append(gachaResultResponse, models.GachaResultResponse{CharacterID: strconv.Itoa(int(characterID)), Name: character.Name})

	c.JSON(http.StatusOK, gin.H{
		"results": gachaResultResponse,
	})

}

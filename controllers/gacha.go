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
	var gachaResultResponse = make([]models.GachaResultResponse, gacha.Times)

	// Get all different types of characters from database
	var ssrCharacters, srCharacters, rCharacters []models.Character

	err = repos.GetAllSpecificCharacters(controller.Db, &ssrCharacters, "ssr")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.GetAllSpecificCharacters(controller.Db, &srCharacters, "sr")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.GetAllSpecificCharacters(controller.Db, &rCharacters, "r")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ssrNumbers, srNumbers, rNumbers := len(ssrCharacters), len(srCharacters), len(rCharacters)

	idToCount := make(map[int]int)

	// Simulate the gacha and record appearance time for each character
	for i := 0; i < int(gacha.Times); i++ {
		n := rand.Intn(100) + 1
		if n == 1 {
			idToCount[int(ssrCharacters[rand.Intn(ssrNumbers)].ID)]++

		} else if n <= 11 {
			idToCount[int(srCharacters[rand.Intn(srNumbers)].ID)]++
		} else {
			idToCount[int(rCharacters[rand.Intn(rNumbers)].ID)]++
		}
	}

	index := 0
	for id, count := range idToCount {
		userCharacters := make([]models.UserCharacter, count)
		var character models.Character
		err := repos.GetCharacter(controller.Db, &character, uint(id))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		idString := strconv.Itoa(id)
		name := character.Name
		for i := range userCharacters {
			userCharacters[i] = models.UserCharacter{UserID: user.ID, CharacterID: uint(id)}
			gachaResultResponse[index] = models.GachaResultResponse{CharacterID: idString, Name: name}
			index++
		}

		err = repos.CreateUserCharacters(controller.Db, &userCharacters)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"results": gachaResultResponse,
	})

}

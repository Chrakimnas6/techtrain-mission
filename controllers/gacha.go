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

	// Get combined information of users
	var charactersOddsComb []struct {
		models.GachaCharacterOdds
		models.Character
	}
	err = repos.GetCharactersOddsComb(controller.Db, &charactersOddsComb)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var gachaResults []models.GachaResult
	var userCharacters []models.UserCharacter

	// Suppose all cards' possibilities sum up to 1 in a certain type of gacha pool
	for i := 0; i < int(gacha.Times); i++ {
		num := rand.Float64()
		oddsSum := 0.0
		for _, character := range charactersOddsComb {
			oddsSum += character.Odds
			if num <= oddsSum {
				userCharacters = append(userCharacters, models.UserCharacter{
					UserID:      user.ID,
					CharacterID: character.CharacterID,
				})
				gachaResults = append(gachaResults, models.GachaResult{
					CharacterID: strconv.Itoa(int(character.CharacterID)),
					Name:        character.Name,
				})
				break
			}
		}
	}

	err = repos.CreateUserCharacters(controller.Db, &userCharacters)

	c.JSON(http.StatusOK, gin.H{
		"results": gachaResults,
	})

}

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
	// Default value
	gacha.ID = 1
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
	err = repos.GetCharactersOddsComb(controller.Db, &charactersOddsComb, gacha.ID)
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

// Create new gacha pool
func (controller *Controller) CreateGachaPool(c *gin.Context) {
	var gachaCharacterOddsAll []models.GachaCharacterOdds
	err := repos.GetDefaultGachaPool(controller.Db, &gachaCharacterOddsAll)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var idNew uint
	err = repos.GetNewestGachaID(controller.Db, &idNew)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Increment id by 1, indicating the new gacha pool
	idNew += 1
	for i, gachaCharacterOdds := range gachaCharacterOddsAll {
		gachaCharacterOdds.GachaID = idNew
		gachaCharacterOddsAll[i] = gachaCharacterOdds
	}
	err = repos.CreateGachaCharacterOddsAll(controller.Db, &gachaCharacterOddsAll)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gachaCharacterOddsAll)
}

// Update possibility
func (controller *Controller) UpdatePossibility(c *gin.Context) {
	var gachaCharacterOdds models.GachaCharacterOdds
	err := c.BindJSON(&gachaCharacterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	err = repos.UpdatePossibility(controller.Db, &gachaCharacterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gachaCharacterOdds)
}

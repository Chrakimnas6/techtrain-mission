package controllers

import (
	"net/http"
	"strconv"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
)

// Get all user_characters in the database
func (controller *Controller) GetUserCharacters(c *gin.Context) {
	// Get all user_characters
	var userCharacters []models.UserCharacter
	err := repos.GetUserCharacters(controller.Db, &userCharacters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var userCharactersResponses []models.UserCharacterResponse
	for _, userCharacter := range userCharacters {
		var character models.Character
		err = repos.GetCharacter(controller.Db, &character, uint(userCharacter.CharacterID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		userCharactersResponses = append(userCharactersResponses, models.UserCharacterResponse{
			UserCharacterID: strconv.Itoa(int(userCharacter.ID)),
			CharacterID:     strconv.Itoa(int(userCharacter.CharacterID)),
			Name:            character.Name,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

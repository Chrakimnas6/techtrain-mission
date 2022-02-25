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

//Get user characters with specfic user's token
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
	var userCharacters []models.UserCharacter
	err = repos.GetUserCharacters(controller.Db, &userCharacters, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// TODO: too time-consuming
	// convert user characters to the format that reponse needs
	var userCharactersResponses = make([]models.UserCharacterResponse, len(userCharacters))
	for index, userCharacter := range userCharacters {
		var character models.Character
		err := repos.GetCharacter(controller.Db, &character, uint(userCharacter.CharacterID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		userCharactersResponses[index] = models.UserCharacterResponse{
			UserCharacterID: strconv.Itoa(int(userCharacter.ID)),
			CharacterID:     strconv.Itoa(int(userCharacter.CharacterID)),
			Name:            character.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

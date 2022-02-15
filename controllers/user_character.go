package controllers

import (
	"errors"
	"net/http"

	//"strconv"

	//"training/helpers"
	"training/helpers"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get all user_characters in the database
func (controller *Controller) GetAllUserCharacters(c *gin.Context) {
	// Get all user_characters
	var userCharacters []models.UserCharacter
	err := repos.GetAllUserCharacters(controller.Db, &userCharacters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var userCharactersResponses []models.UserCharacterResponse
	// for _, userCharacter := range userCharacters {
	// 	// Get character by character ID
	// 	var character models.Character
	// 	err = repos.GetCharacter(controller.Db, &character, uint(userCharacter.CharacterID))
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	// 		return
	// 	}
	// 	// Convert fileds into string according to the reponse
	// 	userCharactersResponses = append(userCharactersResponses, models.UserCharacterResponse{
	// 		UserCharacterID: strconv.Itoa(int(userCharacter.ID)),
	// 		CharacterID:     strconv.Itoa(int(userCharacter.CharacterID)),
	// 		Name:            character.Name,
	// 	})
	// }
	err = helpers.Convert(controller.Db, &userCharacters, &userCharactersResponses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

//Get user_characters with specfic user's token
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
	var userCharacters []models.UserCharacter
	err = repos.GetUserCharacters(controller.Db, &userCharacters, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var userCharactersResponses []models.UserCharacterResponse
	err = helpers.Convert(controller.Db, &userCharacters, &userCharactersResponses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = helpers.Convert(controller.Db, &userCharacters, &userCharactersResponses)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{
		"characters": userCharactersResponses,
	})

}

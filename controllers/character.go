package controllers

import (
	"net/http"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
)

// Get all characters in the database
func (controller *Controller) GetCharacters(c *gin.Context) {
	var characters []models.Character
	err := repos.GetCharacters(controller.Db, &characters)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"characters": characters,
	})
}

// Create a character to the database
func (controller *Controller) CreateCharacter(c *gin.Context) {
	var characterRequest models.CharacterRequest
	err := c.BindJSON(&characterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	character := models.Character{
		Name: characterRequest.Name,
		Rank: characterRequest.Rank,
	}

	err = repos.CreateCharacter(controller.Db, &character)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Create character's odd
	// TODO: Currently only one type of gacha pool and the default id is 1
	characterOdds := models.GachaCharacterOdds{
		GachaID:     1,
		CharacterID: character.ID,
		Odds:        characterRequest.Odds,
	}
	err = repos.CreateGachaCharacterOdds(controller.Db, &characterOdds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, characterOdds)

}

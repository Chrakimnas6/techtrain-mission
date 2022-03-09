package controllers

import (
	"net/http"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
)

type CharacterRequest struct {
	ID   uint    `gorm:"autoIncrement"`
	Name string  `gorm:"not null" json:"name"`
	Rank string  `gorm:"not null" json:"rank"`
	Odds float64 `gorm:"not null" json:"odds"`
}

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
	var characterRequest CharacterRequest
	err := c.BindJSON(&characterRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	character := models.Character{
		Name:          characterRequest.Name,
		CharacterRank: characterRequest.Rank,
	}

	err = repos.CreateCharacter(controller.Db, &character)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	// Create character's odd
	// The default gacha id is 1
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

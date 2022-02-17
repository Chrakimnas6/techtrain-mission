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
	var character models.Character
	err := c.BindJSON(&character)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = repos.CreateCharacter(controller.Db, &character)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, character)

}

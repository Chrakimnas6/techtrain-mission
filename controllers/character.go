package controllers

import (
	"net/http"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
)

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

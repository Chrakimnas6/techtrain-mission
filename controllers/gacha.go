package controllers

import (
	"errors"
	"net/http"
	"training/models"
	"training/repos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Draw with times get from request's body
func (controller *Controller) Draw(c *gin.Context) {
	// Get user with the token
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

	var gacha models.Gacha
	err = c.BindJSON(&gacha)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

}

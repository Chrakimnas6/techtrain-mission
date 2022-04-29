package controllers

import (
	"crypto-tool/transactions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) ConfirmTransaction(c *gin.Context) {
	txHex := c.Query("txHex")
	err := transactions.CheckIfInBlockchain(controller.Client, txHex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Transaction not valid",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Transaction is mined!",
	})

}

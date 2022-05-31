package controllers

import (
	"crypto-tool/models"
	"crypto-tool/repos"
	"crypto-tool/transactions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) ConfirmTransaction(c *gin.Context) {
	// Get the transaction hash from the request body
	txHex := c.Query("txHex")
	// Save the transaction to the database
	var transaction models.Transaction
	transaction.Tx = txHex
	transaction.Status = "pending"
	err := repos.SaveTransaction(controller.Db, transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error saving transaction",
		})
		return
	}
	fmt.Println("Transaction saved!")

	go func() {
		// Check if the transaction is in the blockchain
		fmt.Println("Checking if transaction is in blockchain...")
		err = transactions.CheckIfInBlockchain(controller.Client, txHex)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Transaction not valid",
			})
			return
		}
		fmt.Println("Transaction is mined!")
		// Then change the status in the database to "confirmed"
		err = repos.ChangeTransactionStatus(controller.Db, txHex, "confirmed")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error changing transaction status",
			})
			return
		}
		fmt.Println("Transaction status changed to confirmed!")
	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction sent! Can check status at transaction table in database",
		"tx":      txHex,
	})
}

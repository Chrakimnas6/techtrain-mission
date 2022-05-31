package repos

import (
	"crypto-tool/models"

	"gorm.io/gorm"
)

// Save the transaction
func SaveTransaction(db *gorm.DB, transaction models.Transaction) (err error) {
	err = db.Create(&transaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Change the status of the transaction
func ChangeTransactionStatus(db *gorm.DB, tx string, status string) (err error) {
	var transaction models.Transaction
	err = db.Where("tx = ?", tx).First(&transaction).Error
	if err != nil {
		return err
	}
	transaction.Status = status
	err = db.Save(&transaction).Error
	if err != nil {
		return err
	}
	return nil
}

package repos

import (
	"training/models"

	"gorm.io/gorm"
)

func CheckTransactionIsConfirmed(db *gorm.DB, transaction *models.Transaction, txHex string) (err error) {
	err = db.Where("tx = ? AND status = ?", txHex, "confirmed").First(&transaction).Error
	if err != nil {
		return err
	}
	return nil
}

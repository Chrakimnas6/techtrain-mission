package repos

import (
	"training/models"

	"gorm.io/gorm"
)

func SaveToken(db *gorm.DB, token *models.Token) (err error) {
	err = db.Create(&token).Error
	if err != nil {
		return err
	}
	return nil
}

func GetToken(db *gorm.DB, token *models.Token, symbol string) (err error) {
	err = db.Where("symbol = ?", symbol).Last(&token).Error
	if err != nil {
		return err
	}
	return nil
}

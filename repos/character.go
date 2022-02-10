package repos

import (
	"training/models"

	"gorm.io/gorm"
)

// Create a character
func CreateCharacter(db *gorm.DB, character *models.Character) (err error) {
	err = db.Create(&character).Error
	if err != nil {
		return err
	}
	return nil
}

package repos

import (
	"training/models"

	"gorm.io/gorm"
)

// Create a user_character
func CreateUserCharacter(db *gorm.DB, userCharacter *models.UserCharacter) (err error) {
	err = db.Create(&userCharacter).Error
	if err != nil {
		return err
	}
	return nil
}

// Get all user_characters
func GetUserCharacters(db *gorm.DB, userCharacters *[]models.UserCharacter) (err error) {
	err = db.Find(&userCharacters).Error
	if err != nil {
		return err
	}
	return nil
}

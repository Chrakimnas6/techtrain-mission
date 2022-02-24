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

// Create user_characters
func CreateUserCharacters(db *gorm.DB, userCharacters *[]models.UserCharacter) (err error) {
	err = db.CreateInBatches(&userCharacters, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

// Get all user_characters
func GetAllUserCharacters(db *gorm.DB, userCharacters *[]models.UserCharacter) (err error) {
	err = db.Find(&userCharacters).Error
	if err != nil {
		return err
	}
	return nil
}

// Get specific user's user_characters
func GetUserCharacters(db *gorm.DB, userCharacters *[]models.UserCharacter, userID uint) (err error) {
	err = db.Where("user_id = ?", userID).Find(&userCharacters).Error
	if err != nil {
		return err
	}
	return nil
}

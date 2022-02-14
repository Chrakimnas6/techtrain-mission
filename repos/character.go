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

// Get all users
func GetCharacters(db *gorm.DB, characters *[]models.Character) (err error) {
	err = db.Find(&characters).Error
	if err != nil {
		return err
	}
	return nil
}

// Get table size
func GetSize(db *gorm.DB, size *int64) (err error) {
	err = db.Model(&models.Character{}).Count(size).Error
	if err != nil {
		return err
	}
	return nil

}

// Get character by ID
func GetCharacter(db *gorm.DB, character *models.Character, id uint) (err error) {
	err = db.Model(&models.Character{}).First(&character, id).Error
	if err != nil {
		return err
	}
	return nil
}

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

// Get all characters
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
	//err = db.First(&character, id).Error
	if err != nil {
		return err
	}
	return nil
}

// Get specific type characters
func GetAllSpecificCharacters(db *gorm.DB, characters *[]models.Character, characterType string) (err error) {
	//TODO: why this not work???
	//err = db.Where("rank = ?", characterType).Find(&characters).Error
	err = db.Raw("SELECT * FROM `go_database`.`characters` WHERE `rank` = ?", characterType).Scan(&characters).Error
	if err != nil {
		return err
	}
	return nil
}

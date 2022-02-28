package repos

import (
	"fmt"
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
	//err = db.CreateInBatches(&userCharacters, 1000).Error
	err = db.Create(&userCharacters).Error
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

// Get specific user's user_characters with required response
func GetUserCharacters(db *gorm.DB, userCharactersResponses *[]models.UserCharacterResponse, userID uint) (err error) {
	fmt.Println(userID)
	err = db.Model(&models.UserCharacter{}).Select("user_characters.id as UserCharacterID, characters.id as CharacterID, characters.name as Name").Joins("inner join characters on user_characters.character_id = characters.id").Where("user_characters.user_id = ?", userID).Scan(&userCharactersResponses).Error
	if err != nil {
		return err
	}
	return nil
}

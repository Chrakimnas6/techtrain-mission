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

// Create a character odds
func CreateGachaCharacterOdds(db *gorm.DB, characterOdds *models.GachaCharacterOdds) (err error) {
	err = db.Create(&characterOdds).Error
	if err != nil {
		return err
	}
	return nil
}

// Create character odds
func CreateGachaCharacterOddsAll(db *gorm.DB, characterOdds *[]models.GachaCharacterOdds) (err error) {
	err = db.Create(&characterOdds).Error
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
	//err = db.Where("rank = ?", characterType).Find(&characters).Error
	err = db.Raw("SELECT * FROM `go_database`.`characters` WHERE `rank` = ?", characterType).Scan(&characters).Error
	if err != nil {
		return err
	}
	return nil
}

// Get join information from characters and characters_odds
// TODO: add Where gacha_id = ? for selecting the gacha type
func GetCharactersOddsComb(db *gorm.DB, charactersOddsComb *[]struct {
	models.GachaCharacterOdds
	models.Character
}, gachaID uint) (err error) {
	err = db.Model(&models.GachaCharacterOdds{}).Select("*").Joins("inner join characters on gacha_character_odds.character_id = characters.id").Where("gacha_character_odds.gacha_id = ?", gachaID).Scan(&charactersOddsComb).Error
	if err != nil {
		return err
	}
	return nil

}

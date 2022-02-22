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
	// Also add character to different databases depends on its rank
	switch character.Rank {
	case "ssr":
		character_ssr := &models.CharacterSSR{CharacterID: character.ID}
		err = db.Create(&character_ssr).Error
	case "sr":
		character_sr := &models.CharacterSR{CharacterID: character.ID}
		err = db.Create(&character_sr).Error
	case "r":
		character_r := &models.CharacterR{CharacterID: character.ID}
		err = db.Create(&character_r).Error
	}
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
	if err != nil {
		return err
	}
	return nil
}

// Get random character_id from requested database
func GetRandCharacter(db *gorm.DB, characterID *uint, dbName string) (err error) {
	switch dbName {
	case "ssr":
		err = db.Raw("SELECT character_id FROM character_ssrs ORDER BY RAND() LIMIT 1").Scan(characterID).Error
	case "sr":
		err = db.Raw("SELECT character_id FROM character_srs ORDER BY RAND() LIMIT 1").Scan(characterID).Error
	case "r":
		err = db.Raw("SELECT character_id FROM character_rs ORDER BY RAND() LIMIT 1").Scan(characterID).Error
	}

	if err != nil {
		return err
	}
	return nil
}

// Get how many SSR cards are there in the database
func GetSSRSize(db *gorm.DB, size *int64) (err error) {
	err = db.Model(&models.CharacterSSR{}).Count(size).Error
	if err != nil {
		return err
	}
	return nil
}

// Get how many SSR cards are there in the database
func GetSRSize(db *gorm.DB, size *int64) (err error) {
	err = db.Model(&models.CharacterSR{}).Count(size).Error
	if err != nil {
		return err
	}
	return nil
}

// Get how many SSR cards are there in the database
func GetRSize(db *gorm.DB, size *int64) (err error) {
	err = db.Model(&models.CharacterR{}).Count(size).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSSR(db *gorm.DB, character *models.CharacterSSR, id uint) (err error) {
	err = db.First(&character, id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSR(db *gorm.DB, character *models.CharacterSR, id uint) (err error) {
	err = db.First(&character, id).Error
	if err != nil {
		return err
	}
	return nil
}

func GetR(db *gorm.DB, character *models.CharacterR, id uint) (err error) {
	err = db.First(&character, id).Error
	if err != nil {
		return err
	}
	return nil
}

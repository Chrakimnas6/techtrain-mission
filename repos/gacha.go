package repos

import (
	"training/models"

	"gorm.io/gorm"
)

// Update odds
func UpdateOdds(db *gorm.DB, gachaCharacterOdds *models.GachaCharacterOdds) (err error) {
	err = db.Model(&models.GachaCharacterOdds{}).
		Where("gacha_id = ? AND character_id = ?", gachaCharacterOdds.GachaID, gachaCharacterOdds.CharacterID).
		Update("odds", gachaCharacterOdds.Odds).Error
	if err != nil {
		return err
	}
	return nil
}

// Create new gacha pool
// 1. Get the default gacha pool
func GetDefaultGachaPool(db *gorm.DB, gachaCharacterOdds *[]models.GachaCharacterOdds) (err error) {
	// The default gacha pool's id is 1
	err = db.Model(&models.GachaCharacterOdds{}).Where("gacha_id = ?", 1).Find(&gachaCharacterOdds).Error
	if err != nil {
		return err
	}
	return nil
}

// 2. Find the newest gacha id in the database
func GetNewestGachaID(db *gorm.DB, idNew *uint) (err error) {
	err = db.Model(&models.GachaCharacterOdds{}).Select("max(gacha_id)").Row().Scan(idNew)
	if err != nil {
		return err
	}
	return nil
}

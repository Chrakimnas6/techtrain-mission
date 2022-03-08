package controllers

import (
	// "training/models"

	"gorm.io/gorm"
)

type Controller struct {
	Db *gorm.DB
}

// Initialize new database of User
func New(db *gorm.DB) *Controller {
	// // Create user table
	// db.AutoMigrate(&models.User{})
	// // Create character table
	// db.AutoMigrate(&models.Character{})
	// // Create user_character table
	// db.AutoMigrate(&models.UserCharacter{})
	// // Create gacha_character_odds table
	// db.AutoMigrate(&models.GachaCharacterOdds{})
	return &Controller{Db: db}
}

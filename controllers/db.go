package controllers

import (
	"training/models"

	"gorm.io/gorm"
)

type Controller struct {
	Db *gorm.DB
}

// Initialize new database of User
func New(db *gorm.DB) *Controller {
	// Create user table
	db.AutoMigrate(&models.User{})
	// Create character table
	db.AutoMigrate(&models.Character{})
	// Create user_character table
	db.AutoMigrate(&models.UserCharacter{})
	// Create tables for different types of cards
	db.AutoMigrate(&models.CharacterSSR{})
	db.AutoMigrate(&models.CharacterSR{})
	db.AutoMigrate(&models.CharacterR{})
	return &Controller{Db: db}
}

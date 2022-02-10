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
	return &Controller{Db: db}
}

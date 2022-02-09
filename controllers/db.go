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
	return &Controller{Db: db}
}

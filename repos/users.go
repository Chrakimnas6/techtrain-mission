package repos

import (
	"training/models"

	"gorm.io/gorm"
)

// Create a user
func CreateUser(db *gorm.DB, user *models.User) (err error) {
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Get all users
func GetUsers(db *gorm.DB, users *[]models.User) (err error) {
	err = db.Find(&users).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by token
func GetUser(db *gorm.DB, user *models.User, token string) (err error) {
	err = db.Where("token = ?", token).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Get user by name
func GetUserByName(db *gorm.DB, user *models.User, name string) (err error) {
	err = db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Update user
func UpdateUser(db *gorm.DB, user *models.User) (err error) {
	err = db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

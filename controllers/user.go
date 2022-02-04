package controllers

import (
	"errors"
	"fmt"
	"training/db"
	"training/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	Db *gorm.DB
}

// Initialize new database of User
func New() *UserRepo {
	db := db.Init()
	db.AutoMigrate(&models.User{})
	return &UserRepo{Db: db}
}

// Get all users in the database
func (repository *UserRepo) GetUsers(c *gin.Context) {
	var users []models.User
	err := models.GetUsers(repository.Db, &users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.HTML(200, "index.html", gin.H{
		"users": users,
	})
}

// Create user with "name"
// Use uuid to generate token for the user
func (repository *UserRepo) CreateUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	token := uuid.New()
	user.Token = token.String()
	fmt.Println("create user " + user.Name + " with token " + user.Token)
	err := models.CreateUser(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": user.Token,
	})
}

// Get user by the token from Header
func (repository *UserRepo) GetUser(c *gin.Context) {
	var user models.User
	token := c.GetHeader("x-token")
	err := models.GetUser(repository.Db, &user, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name": user.Name,
	})
}

// Update user by the token from Header with new name
func (repository *UserRepo) UpdateUser(c *gin.Context) {
	var user models.User
	token := c.GetHeader("x-token")
	err := models.GetUser(repository.Db, &user, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&user)
	err = models.UpdateUser(repository.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Status(http.StatusOK)
}

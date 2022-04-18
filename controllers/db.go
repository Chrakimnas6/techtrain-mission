package controllers

import (
	// "training/models"
	token "training/contracts"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type Controller struct {
	Db       *gorm.DB
	Client   *ethclient.Client
	Keystore *keystore.KeyStore
	Instance *token.Token
}

// func New(db *gorm.DB) *Controller {
// 	// Create user table
// 	db.AutoMigrate(&models.User{})
// 	// Create character table
// 	db.AutoMigrate(&models.Character{})
// 	// Create user_character table
// 	db.AutoMigrate(&models.UserCharacter{})
// 	// Create gacha_character_odds table
// 	db.AutoMigrate(&models.GachaCharacterOdds{})
// 	return &Controller{Db: db}
// }

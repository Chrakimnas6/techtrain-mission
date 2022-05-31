package controllers

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

type Controller struct {
	Db     *gorm.DB
	Client *ethclient.Client
}

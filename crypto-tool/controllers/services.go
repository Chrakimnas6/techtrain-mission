package controllers

import "github.com/ethereum/go-ethereum/ethclient"

type Controller struct {
	Client *ethclient.Client
}

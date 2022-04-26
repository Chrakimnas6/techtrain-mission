package token

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckInformation(instance *Token) {
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	supply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Token's name: %s\n", name)
	fmt.Printf("Token's symbol: %s\n", symbol)
	fmt.Printf("Token's supply: %s\n", supply)

}

// Get address' ETH balance
func GetETHBalance(client *ethclient.Client, address common.Address) (balance *big.Int) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	return balance
}

// Get address' token balance
func GetTokenBalance(instance *Token, address common.Address) (balance *big.Int) {
	balance, err := instance.BalanceOf(nil, address)
	if err != nil {
		log.Fatal(err)
	}
	return new(big.Int).Div(balance, big.NewInt(1000000000000000000))
}

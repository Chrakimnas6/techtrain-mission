package token

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func CheckInformation(instance *Token, address common.Address) {
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

	balanceA, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Token's name: %s\n", name)
	fmt.Printf("Token's symbol: %s\n", symbol)
	fmt.Printf("Token's supply: %s\n", supply)

	fmt.Printf("A's Balance is: %s\n", new(big.Int).Div(balanceA, big.NewInt(1000000000000000000)))
}

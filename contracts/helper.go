package token

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckInformation(instance *Token) (name string, symbol string, supply *big.Int, err error) {
	name, err = instance.Name(&bind.CallOpts{})
	if err != nil {
		return "", "", nil, err
	}

	symbol, err = instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return "", "", nil, err
	}
	supply, err = instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return "", "", nil, err
	}

	fmt.Printf("Token's name: %s\n", name)
	fmt.Printf("Token's symbol: %s\n", symbol)
	fmt.Printf("Token's supply: %s\n", supply)

	return name, symbol, supply, err
}

// Get address' ETH balance
func GetETHBalance(client *ethclient.Client, address common.Address) (balance *big.Int, err error) {
	balance, err = client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Get address' token balance
func GetTokenBalance(instance *Token, address common.Address) (balance *big.Int, err error) {
	balance, err = instance.BalanceOf(nil, address)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// Check if the transaction is mined
func CheckTransaction(txHash *types.Transaction) (err error) {
	response, err := http.Get("http://crypto-tool:8888/transaction/confirm?txHex=" + txHash.Hash().Hex())
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("transaction not valid")
	}
}

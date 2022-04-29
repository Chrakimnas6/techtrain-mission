package token

import (
	"context"
	"fmt"
	"math/big"
	"training/accounts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Deploy(client *ethclient.Client, ks *keystore.KeyStore, keystoreFileName string) (tx *types.Transaction, tokenAddress common.Address, instance *Token, err error) {
	fmt.Println("Importing account...")
	account, err := accounts.ImportAccount(ks, keystoreFileName, "password")
	_ = err
	// if err != nil {
	// 	fmt.Println("Error importing account:", err)
	// 	return nil, tokenAddress, nil, err
	// }
	ks.Unlock(account, "password")
	fmt.Println("Account imported")

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, tokenAddress, nil, err
	}
	fmt.Printf("nonce: %d\n", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, tokenAddress, nil, err
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, tokenAddress, nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, tokenAddress, nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	fmt.Println("Deploying contract...")
	tokenAddress, tx, instance, err = DeployToken(auth, client)
	if err != nil {
		return nil, tokenAddress, nil, err
	}
	fmt.Printf("Token address is: %s\n", tokenAddress.Hex())
	return tx, tokenAddress, instance, nil
}

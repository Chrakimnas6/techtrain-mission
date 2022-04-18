package token

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"training/accounts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Deploy(client *ethclient.Client, ks *keystore.KeyStore, keystoreFileName string) (tokenAddress common.Address, instance *Token) {

	account := accounts.ImportAccount(ks, keystoreFileName, "password")

	ks.Unlock(account, "password")
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d\n", nonce)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	tokenAddress, tx, instance, err := DeployToken(auth, client)
	if err != nil {
		log.Fatal(err)
	}
	_ = tx
	fmt.Printf("Token address is: %s\n", tokenAddress.Hex())
	return tokenAddress, instance
}

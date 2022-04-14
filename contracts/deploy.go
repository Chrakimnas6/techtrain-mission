package token

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"training/accounts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Deploy(keystoreFileName string) (tokenAddress common.Address, instance *Token) {
	client, err := ethclient.Dial("http://hardhat:8545/")
	//nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		log.Fatal(err)
	}

	ks, account := accounts.ImportKeystore(keystoreFileName, "password")

	ks.Unlock(account, "password")
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
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

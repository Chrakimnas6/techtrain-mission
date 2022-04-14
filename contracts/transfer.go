package token

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"training/accounts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func FaucetTransfer(adminAddress common.Address) (err error) {

	client, err := ethclient.Dial("http://hardhat:8545/")
	if err != nil {
		return err
	}
	fmt.Println("Connected to Ethereum client")
	value := big.NewInt(1000000000000000000)
	privateKeyFrom, addressFrom := accounts.GetHardhatAddress()

	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
	if err != nil {
		return err
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, adminAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction with the private key of the sender
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	//
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyFrom)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}

func TransferETH(privateKeyFrom *ecdsa.PrivateKey, addressFrom common.Address,
	addressTo common.Address, value *big.Int, client *ethclient.Client) (err error) {

	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
	if err != nil {
		return err
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, addressTo, value, gasLimit, gasPrice, nil)

	// Sign the transaction with the private key of the sender
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return err
	}
	//
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyFrom)
	if err != nil {
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}

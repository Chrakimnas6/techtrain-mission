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

// Receive ETH from Hardhat's address
// func FaucetTransfer(client *ethclient.Client, adminAddress common.Address) (err error) {
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Connected to Ethereum client")
// 	value := big.NewInt(1000000000000000000)
// 	privateKeyFrom, addressFrom := accounts.GetHardhatAddress()

// 	nonce, err := client.PendingNonceAt(context.Background(), addressFrom)
// 	if err != nil {
// 		return err
// 	}
// 	gasLimit := uint64(21000)
// 	gasPrice, err := client.SuggestGasPrice(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	tx := types.NewTransaction(nonce, adminAddress, value, gasLimit, gasPrice, nil)

// 	// Sign the transaction with the private key of the sender
// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyFrom)
// 	if err != nil {
// 		return err
// 	}

// 	err = client.SendTransaction(context.Background(), signedTx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func TransferETH(client *ethclient.Client, ks *keystore.KeyStore, keystoreFileName string, addressTo common.Address, amount int) (tx *types.Transaction, err error) {
	account, err := accounts.ImportAccount(ks, keystoreFileName, "password")
	_ = err
	// if err != nil {
	// 	return nil, err
	// }
	ks.Unlock(account, "password")

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Amount: %d\n", amount)
	// multiply by 10 * 16
	transferAmount := big.NewInt(0).Mul(big.NewInt(int64(amount)), big.NewInt(10000000000000000))
	fmt.Printf("Transfer amount: %d\n", transferAmount)

	tx = types.NewTransaction(nonce, addressTo, transferAmount, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	// Sign the transaction using keystore
	signedTx, err := ks.SignTx(account, tx, chainID)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

// Transfer token from admin
func TransferToken(client *ethclient.Client, ks *keystore.KeyStore, instance *Token, keystoreFileName string, addressTo common.Address, amount int) (tx *types.Transaction, err error) {
	account, err := accounts.ImportAccount(ks, keystoreFileName, "password")
	_ = err
	// if err != nil {
	// 	return nil, err
	// }
	ks.Unlock(account, "password")

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	fmt.Printf("Amount: %d\n", amount)
	transferAmount := big.NewInt(0).Mul(big.NewInt(int64(amount)), big.NewInt(1000000000000000000))
	fmt.Printf("Transfer amount: %d\n", transferAmount)

	tx, err = instance.Transfer(auth, addressTo, transferAmount)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Consume token when executing the gacha
func BurnToken(client *ethclient.Client, ks *keystore.KeyStore, instance *Token, keystoreFileName string, amount int) (tx *types.Transaction, err error) {
	account, err := accounts.ImportAccount(ks, keystoreFileName, "password")
	_ = err
	// if err != nil {
	// 	return nil, err
	// }
	ks.Unlock(account, "password")

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	fmt.Printf("Amount: %d\n", amount)
	burnAmount := big.NewInt(0).Mul(big.NewInt(int64(amount)), big.NewInt(1000000000000000000))
	fmt.Printf("Transfer amount: %d\n", burnAmount)

	tx, err = instance.Burn(auth, burnAmount)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// Mint Token
func MintToken(client *ethclient.Client, ks *keystore.KeyStore, instance *Token, keystoreFileName string, amount int) (tx *types.Transaction, err error) {
	account, err := accounts.ImportAccount(ks, keystoreFileName, "password")
	_ = err
	// if err != nil {
	// 	return nil, err
	// }
	ks.Unlock(account, "password")

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyStoreTransactorWithChainID(ks, account, chainID)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	fmt.Printf("Amount: %d\n", amount)
	mintAmount := big.NewInt(0).Mul(big.NewInt(int64(amount)), big.NewInt(1000000000000000000))
	fmt.Printf("Mint amount: %d\n", mintAmount)

	tx, err = instance.Mint(auth, account.Address, mintAmount)
	if err != nil {
		return nil, err
	}
	return tx, err
}

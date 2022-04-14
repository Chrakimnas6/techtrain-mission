package accounts

import (
	"context"
	"crypto/ecdsa"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Get account from Hardhat local node
func GetHardhatAddress() (privateKey *ecdsa.PrivateKey, address common.Address) {
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}
	address = deriveAddress(privateKey)
	return privateKey, address
}

// Derive address from private key
func deriveAddress(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address
}

// Create new keystore
func CreateKeystore(password string) (ks *keystore.KeyStore, keystoreFileName string, acount accounts.Account) {
	ks = keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	keystoreFileName = account.URL.Path
	return ks, keystoreFileName, account
}

// Import keystore
func ImportKeystore(keystoreFileName string, password string) (ks *keystore.KeyStore, account accounts.Account) {
	ks = keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(keystoreFileName)
	if err != nil {
		log.Fatal(err)
	}
	account, err = ks.Import(jsonBytes, password, password)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//address = account.Address
	return ks, account
}

// Get address' balance
func GetBalance(address common.Address) (balance *big.Int) {
	client, err := ethclient.Dial("http://hardhat:8545/")
	if err != nil {
		log.Fatal(err)
	}
	balance, err = client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	return balance
}

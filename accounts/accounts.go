package accounts

import (
	"crypto/ecdsa"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

// Create new account
func CreateAccount(ks *keystore.KeyStore, password string) (keystoreFileName string, acount accounts.Account) {
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}
	keystoreFileName = account.URL.Path
	return keystoreFileName, account
}

// Get account
func ImportAccount(ks *keystore.KeyStore, keystoreFileName string, password string) (account accounts.Account) {
	jsonBytes, err := ioutil.ReadFile(keystoreFileName)
	if err != nil {
		log.Fatal(err)
	}
	account, err = ks.Import(jsonBytes, password, password)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return account
}

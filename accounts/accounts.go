package accounts

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// Get account from Hardhat local node
// func GetHardhatAddress() (privateKey *ecdsa.PrivateKey, address common.Address) {
// 	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	address = deriveAddress(privateKey)
// 	return privateKey, address
// }

// Derive address from private key
// func deriveAddress(privateKey *ecdsa.PrivateKey) common.Address {
// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Fatal("error casting public key to ECDSA")
// 	}
// 	address := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	return address
// }

// Create new account
func CreateAccount(ks *keystore.KeyStore, password string) (keystoreFileName string, acount accounts.Account, err error) {
	account, err := ks.NewAccount(password)
	if err != nil {
		return "", account, err
	}
	keystoreFileName = account.URL.Path
	return keystoreFileName, account, nil
}

// Get account
func ImportAccount(ks *keystore.KeyStore, keystoreFileName string, password string) (account accounts.Account, err error) {
	jsonBytes, err := ioutil.ReadFile(keystoreFileName)
	if err != nil {
		return account, err
	}
	account, err = ks.Import(jsonBytes, password, password)
	_ = err
	return account, err
}

// Stores existing admin's key into the key directory
func StoreKey(ks *keystore.KeyStore, password string) (account accounts.Account, err error) {
	privateKey, err := crypto.HexToECDSA("*****")
	if err != nil {
		return account, err
	}
	account, err = ks.ImportECDSA(privateKey, "password")
	if err != nil {
		return account, err
	}
	return account, nil
}

package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckIfInBlockchain(client *ethclient.Client, txHex string) (err error) {
	txHash := common.HexToHash(txHex)
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
			_ = tx
			if err != nil {
				return err
			}
			if !isPending {
				fmt.Println("Transaction mined!")
				return nil
			}
		case <-quit:
			ticker.Stop()
			return err
		}
	}
}

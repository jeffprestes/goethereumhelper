package goethereumhelper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//WaitForTransactionProcessing check trx mining and return his results
func WaitForTransactionProcessing(client *ethclient.Client, trx *types.Transaction, maxAttempts int, interval int) (txReceipt *types.Receipt, err error) {
	var isPending = true
	for isPending {
		fmt.Print(".")
		time.Sleep(time.Duration(interval) * time.Second)
		_, isPending, err = client.TransactionByHash(context.Background(), trx.Hash())
		if err != nil {
			log.Println("[WaitForTransctionProcessing] Error checking if a transaction is mining pending. Error: ", err)
			return
		}
		if !isPending {
			break
		}
		maxAttempts--
		if maxAttempts == 0 {
			log.Println("[WaitForTransctionProcessing] Attempts number exceeded max attempts limit. Error: ", err)
			return
		}
	}
	fmt.Print("\n")
	txReceipt, err = client.TransactionReceipt(context.Background(), trx.Hash())
	if err != nil {
		log.Println("[WaitForTransctionProcessing] It was not possible to get add info category transaction receipt. Error: ", err.Error())
		return
	}
	if txReceipt.Status < 1 {
		err = fmt.Errorf("Transaction failed. Status: %d", txReceipt.Status)
		log.Printf("[WaitForTransctionProcessing] %s\n", err.Error())
		return
	}
	return
}

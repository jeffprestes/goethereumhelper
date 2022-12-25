package goethereumhelper

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WaitForTransactionProcessing check trx mining and return his results
func WaitForTransactionProcessing(client *ethclient.Client, trx *types.Transaction, maxAttempts int, interval int) (txReceipt *types.Receipt, err error) {
	isPending := true
	var ci int
	for isPending {
		var cs string
		if ci == 0 {
			cs = "|"
		} else if ci == 1 {
			cs = "/"
		} else if ci == 2 {
			cs = "-"
		} else if ci == 3 {
			ci = 0
			cs = "|"
		}
		ci++
		if runtime.GOOS == "windows" {
			fmt.Print(".")
		} else {
			fmt.Print("\033[1D" + cs)
		}
		time.Sleep(time.Duration(interval) * time.Second)
		_, isPending, err = client.TransactionByHash(context.Background(), trx.Hash())
		if err != nil {
			err = fmt.Errorf("error getting Tx data by hash: %s", err.Error())
			log.Println("[WaitForTransctionProcessing] Error checking if a transaction is mining pending. Error: ", err)
			return
		}
		if !isPending {
			txReceipt, err = client.TransactionReceipt(context.Background(), trx.Hash())
			if err != nil {
				// In case the node is outdate and have not received the transaction yet
				if !strings.Contains(err.Error(), "not found") {
					err = fmt.Errorf("error getting Tx Receipt: %s", err.Error())
					log.Println("[WaitForTransctionProcessing] It was not possible to get add info category transaction receipt. Error: ", err.Error())
					return
				} else {
					isPending = true
				}
			}
			if !isPending {
				break
			}
		}

		maxAttempts--
		if maxAttempts == 0 {
			err = fmt.Errorf("attempts number exceeded max attempts limit: %d", maxAttempts)
			log.Println("[WaitForTransctionProcessing] Error: ", err)
			return
		}
	}
	fmt.Print("\033[1D")

	if txReceipt.Status < 1 {
		err = fmt.Errorf("transaction in blockchain failed. Status: %d", txReceipt.Status)
		log.Printf("[WaitForTransctionProcessing] %s\n", err.Error())
		return
	}
	return
}

// GetTransactionResult check trx mining and return his results
func GetTransactionResult(client *ethclient.Client, trx common.Hash, maxAttempts int, interval int) (txReceipt *types.Receipt, err error) {
	isPending := true
	var ci int
	for isPending {
		var cs string
		if ci == 0 {
			cs = "|"
		} else if ci == 1 {
			cs = "/"
		} else if ci == 2 {
			cs = "-"
		} else if ci == 3 {
			ci = 0
			cs = "|"
		}
		ci++
		if runtime.GOOS == "windows" {
			fmt.Print(".")
		} else {
			fmt.Print("\033[1D" + cs)
		}
		time.Sleep(time.Duration(interval) * time.Second)
		_, isPending, err = client.TransactionByHash(context.Background(), trx)
		if err != nil && strings.TrimSpace(err.Error()) != "not found" {
			log.Println("[GetTransactionResult] Error checking if a transaction is mining pending. Error: ", err)
			return
		}
		if !isPending {
			break
		}
		maxAttempts--
		if maxAttempts == 0 {
			err = fmt.Errorf("attempts number exceeded max attempts limit: %d", maxAttempts)
			log.Println("[GetTransactionResult] Error maxAttempts: ", err)
			return
		}
	}
	fmt.Print("\033[1D")
	txReceipt, err = client.TransactionReceipt(context.Background(), trx)
	if err != nil {
		log.Println("[GetTransactionResult] It was not possible to get ", trx.String(), "transaction receipt. Error: ", err.Error())
		return
	}
	if txReceipt.Status < 1 {
		err = fmt.Errorf("transaction failed. Status: %d", txReceipt.Status)
		log.Printf("[GetTransactionResult] Status %s\n", err.Error())
		return
	}
	return
}

package goethereumhelper

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//SendEther an example that shows how to send ether from an account to another using Go (Golang)
func SendEther(client *ethclient.Client, senderPrivateKey *ecdsa.PrivateKey, to common.Address, value int64) (receipt *types.Receipt, err error) {

	publicKey := senderPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		return
	}
	// fmt.Println("from", fromAddress.Hex())

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	toAddress := to
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(value), gasLimit, gasPrice, data)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), senderPrivateKey)
	if err != nil {
		log.Println(err)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println(err)
		return
	}

	receipt, err = client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		log.Println(err)
		return
	}
	if receipt == nil {
		err = errors.New("receipt is nil. Forgot to commit?")
		log.Println(err)
		return
	}
	return

	//fmt.Printf("status: %v\n", receipt.Status) // status: 1
}

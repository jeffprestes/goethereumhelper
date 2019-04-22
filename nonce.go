package goethereumhelper

import (
	"context"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//GetNonceNumber gets actual nonce number of an Ethereum address/account
func GetNonceNumber(client *ethclient.Client, pubkey ecdsa.PublicKey) (nonce uint64, err error) {
	err = nil
	origem := crypto.PubkeyToAddress(pubkey)
	nonce, err = client.PendingNonceAt(context.Background(), origem)
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao gerar o nonce da conta na rede: %+v", err)
		return
	}
	return
}

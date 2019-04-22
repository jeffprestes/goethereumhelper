package goethereumhelper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//NewAccount generates new Ethereum address/account
func NewAccount() (address common.Address) {
	key, _ := crypto.GenerateKey()
	address = crypto.PubkeyToAddress(key.PublicKey)
	return
}

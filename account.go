package goethereumhelper

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewAccount generates new Ethereum address/account
func NewAccount() (accountPrivateKey *ecdsa.PrivateKey, address common.Address, err error) {
	accountPrivateKey, err = crypto.GenerateKey()
	if err != nil {
		return nil, address, err
	}
	address = crypto.PubkeyToAddress(accountPrivateKey.PublicKey)
	return
}

// GetPubKey gets public key and address from a private key
func GetPubKey(accountPrivateKey *ecdsa.PrivateKey) (publicKeyECDSA *ecdsa.PublicKey, address common.Address, err error) {
	publicKey := accountPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA)
	return
}

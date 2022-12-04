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

// SendEtherUsingKeystoreWallet an example that shows how to send ether using an account from KeystoreWallet to another using Go (Golang)
func SendEtherUsingKeystoreWallet(client *ethclient.Client, sender KeystoreWallet, to common.Address, value int64) (signedTx *types.Transaction, err error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error getting chainID ", err)
		return
	}

	latestEthBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error gettin the latest Eth Block Header ", err)
		return
	}

	nonce, err := client.PendingNonceAt(context.Background(), sender.Account.Address)
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error getting nonce ", err)
		return
	}

	gasLimit := uint64(21000) // in units
	// Use new EIP-1559
	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error getting Blockchain suggested gasTip ", err)
		return
	}

	maxGasFeeAccepted := new(big.Int).Add(
		latestEthBlockHeader.BaseFee,
		gasTip)

	var data []byte

	eip1559Tx := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		Gas:       gasLimit,
		GasFeeCap: maxGasFeeAccepted,
		GasTipCap: gasTip,
		To:        &to,
		Value:     big.NewInt(value),
		Data:      data,
	}
	tx := types.NewTx(&eip1559Tx)

	signedTx, err = sender.Wallet.SignTx(sender.Account, tx, chainID)
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error Signing Transaction ", err)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("SendEtherUsingKeystoreWallet - Error sending Transaction ", err)
		return
	}
	return
}

func SendEtherUsingPrivateKey(client *ethclient.Client, senderPrivateKey *ecdsa.PrivateKey, to common.Address, value int64) (signedTx *types.Transaction, err error) {
	signedTx, err = SendEtherUsingPrivateKeyGasTipFactor(client, senderPrivateKey, to, 1, value)
	return
}

// SendEtherUsingPrivateKeyGasTipFactor an example that shows how to send ether using private key from an account to another using Go (Golang) with GasTip price factor
func SendEtherUsingPrivateKeyGasTipFactor(client *ethclient.Client, senderPrivateKey *ecdsa.PrivateKey, to common.Address, gasTipFactor int64, value int64) (signedTx *types.Transaction, err error) {
	publicKey := senderPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error getting chainID ", err)
		return
	}

	latestEthBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error gettin the latest Eth Block Header ", err)
		return
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error getting nonce ", err)
		return
	}

	gasLimit := uint64(21000) // in units
	// Use new EIP-1559
	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error getting Blockchain suggested gasTip ", err)
		return
	}

	gasTip = gasTip.Mul(gasTip, big.NewInt(gasTipFactor))

	maxGasFeeAccepted := new(big.Int).Add(
		latestEthBlockHeader.BaseFee,
		gasTip)

	var data []byte

	eip1559Tx := types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		Gas:       gasLimit,
		GasFeeCap: maxGasFeeAccepted,
		GasTipCap: gasTip,
		To:        &to,
		Value:     big.NewInt(value),
		Data:      data,
	}
	tx := types.NewTx(&eip1559Tx)

	signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(chainID), senderPrivateKey)
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error Signing Transaction ", err)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Println("SendEtherUsingPrivateKey - Error sending Transaction ", err)
		return
	}
	return
}

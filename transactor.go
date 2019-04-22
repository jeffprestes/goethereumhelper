package goethereumhelper

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
UpdateKeyedTransactor updates a keyed (signed?) transctor do perform a transaction within a Simulated Ethereum Blockchain
*/
func UpdateKeyedTransactor(transactor *bind.TransactOpts, backend *backends.SimulatedBackend, increaseNonceFactor int, valueToSend int) (err error) {
	err = nil
	nonce, err := backend.PendingNonceAt(context.Background(), transactor.From)
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao obter o nonce para o endereco %s da rede: %+v", transactor.From.String(), err)
		return
	}
	gasPrice, err := backend.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao obter o preco sugerido de gas da rede: %+v", err)
		return
	}

	transactor.GasLimit = uint64(6869310)
	transactor.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	transactor.Value = big.NewInt(int64(valueToSend))
	transactor.Nonce = big.NewInt(int64(nonce))
	return
}

/*
GetKeyedTransactorRinkeby gets a keyed (signed?) transctor do perform a transaction within the Rinkeby Ethereum Blockchain
*/
func GetKeyedTransactorRinkeby(client *ethclient.Client, increaseNonceFactor int) (transactor *bind.TransactOpts, err error) {
	err = nil

	pvtkey, err := crypto.HexToECDSA(os.Getenv("privatekey"))
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao gerar a chave privada: %+v", err)
		return
	}
	pubkey := pvtkey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("[GetKeyedTransactor] Houve falha fazer o casting da chave publica para o padrão ECDSA")
		log.Printf(err.Error())
		return
	}
	basicNonce, err := GetNonceNumber(client, *pubkeyECDSA)
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao obter o nonce da rede: %+v", err)
		return
	}
	nonce := basicNonce + uint64(increaseNonceFactor)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao obter o preco sugerido de gas da rede: %+v", err)
		return
	}

	transactor = bind.NewKeyedTransactor(pvtkey)
	transactor.GasLimit = uint64(6869310)
	transactor.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(2))
	transactor.Value = big.NewInt(0)
	transactor.Nonce = big.NewInt(int64(nonce))

	return
}

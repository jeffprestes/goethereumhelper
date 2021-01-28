package goethereumhelper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func getMockAccountWithFunds(genesisAccount *bind.TransactOpts, backend *backends.SimulatedBackend) (newAccount *bind.TransactOpts, err error) {
	key, _ := crypto.GenerateKey()
	newAccount = bind.NewKeyedTransactor(key)
	err = prepareNewTransction(genesisAccount, backend)
	if err != nil {
		fmt.Printf("Houve falha ao obter novo nonce para fazer a inclusao de saldo na conta de um novo 3rd-party: %+v", err)
		return
	}
	tx := types.NewTransaction(uint64(genesisAccount.Nonce.Int64()), newAccount.From, big.NewInt(101000000000000), genesisAccount.GasLimit, genesisAccount.GasPrice, nil)
	signedTx, err := genesisAccount.Signer(genesisAccount.From, tx)
	if err != nil {
		fmt.Printf("Houve falha ao assinar a transação para enviar ether a nova conta: %+v", err)
		return
	}
	backend.SendTransaction(context.Background(), signedTx)
	backend.Commit()
	return
}

func prepareNewTransction(auth *bind.TransactOpts, backend *backends.SimulatedBackend) (err error) {
	err = nil
	nonce, err := backend.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return
	}
	gasPrice, err := backend.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	auth.GasLimit = uint64(600000)
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	return
}

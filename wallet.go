package goethereumhelper

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// KeystoreWallet implements the accounts.Wallet interface for the original
type KeystoreWallet struct {
	Account  accounts.Account   // Single account contained in this wallet
	Keystore *keystore.KeyStore // Keystore where the account originates from
	Wallet   accounts.Wallet
}

//NewKeystoreWallet returns new instance of KeystoreWallet
func NewKeystoreWallet(ks *keystore.KeyStore, accountHex string) (ksw *KeystoreWallet, err error) {
	ksw = new(KeystoreWallet)
	a1 := accounts.Account{}
	a1.Address = common.HexToAddress(accountHex)
	account, err := ks.Find(a1)
	if err != nil {
		log.Fatalln("NewKeystoreWallet - Conta não existente na keystore", err.Error())
	}
	for _, w := range ks.Wallets() {
		for _, acc := range w.Accounts() {
			if acc.Address.Hash() == account.Address.Hash() {
				ksw.Wallet = w
				break
			}
		}
	}
	ksw.Keystore = ks
	ksw.Account = account
	return
}

// URL implements accounts.Wallet, returning the URL of the account within.
func (w *KeystoreWallet) URL() accounts.URL {
	return w.Account.URL
}

// SignTextWithPassphrase implements accounts.Wallet, attempting to sign the
// given hash with the given account using passphrase as extra authentication.
func (w *KeystoreWallet) SignTextWithPassphrase(passphrase string, text []byte) ([]byte, error) {
	// Account seems valid, request the keystore to sign
	return w.Keystore.SignHashWithPassphrase(w.Account, passphrase, accounts.TextHash(text))
}

// SignTxWithPassphrase implements accounts.Wallet, attempting to sign the given
// transaction with the given account using passphrase as extra authentication.
func (w *KeystoreWallet) SignTxWithPassphrase(passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Account seems valid, request the keystore to sign
	return w.Keystore.SignTxWithPassphrase(w.Account, passphrase, tx, chainID)
}

// SignData signs keccak256(data). The mimetype parameter describes the type of data being signed
func (w *KeystoreWallet) SignData(mimeType string, data []byte) ([]byte, error) {
	return w.Keystore.SignHash(w.Account, crypto.Keccak256(data))
}

//GetNonceNumber gets actual nonce number of an Ethereum address/account
func (w *KeystoreWallet) GetNonceNumber(client *ethclient.Client) (nonce uint64, err error) {
	err = nil
	nonce, err = client.PendingNonceAt(context.Background(), w.Account.Address)
	if err != nil {
		return
	}
	return
}

/*
UpdateKeyedTransactor updates a keyed (signed?) transctor do perform a transaction within a Simulated Ethereum Blockchain
*/
func (w *KeystoreWallet) UpdateKeyedTransactor(transactor *bind.TransactOpts, client *ethclient.Client, increaseNonceFactor int, valueToSend int) (err error) {
	err = nil
	nonce, err := client.PendingNonceAt(context.Background(), transactor.From)
	if err != nil {
		log.Printf("[GetKeyedTransactor] Houve falha ao obter o nonce para o endereco %s da rede: %+v", transactor.From.String(), err)
		return
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
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

// NewKeyStoreTransactor is a utility method to easily create a transaction signer from
// an decrypted key from a keystore
func (w *KeystoreWallet) NewKeyStoreTransactor(passphrase string) (*bind.TransactOpts, error) {
	return &bind.TransactOpts{
		From: w.Account.Address,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != w.Account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := w.Keystore.SignHashWithPassphrase(w.Account, passphrase, signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

//GenerateSignedTxAsJSON generates a JSON signed raw Ethereum Transaction
func (w *KeystoreWallet) GenerateSignedTxAsJSON(
	txOpts *bind.TransactOpts,
	passphrase string,
	contractMethodParameters []byte,
	smartContractAddress common.Address,
	nonce, chainID uint64) (txJSON []byte, err error) {

	tx := types.NewTransaction(nonce, smartContractAddress, big.NewInt(0), txOpts.GasLimit, txOpts.GasPrice, contractMethodParameters)

	txSigned, err := w.SignTxWithPassphrase(passphrase, tx, big.NewInt(int64(chainID)))
	if err != nil {
		log.Fatalln("Erro obter o signer da conta", w.Account.Address.Hex(), err.Error())
		return
	}

	txJSON, err = txSigned.MarshalJSON()
	if err != nil {
		log.Fatalln("Error serializing the transaction", err.Error())
		return
	}
	return
}

package goethereumhelper

import (
	"context"
	"errors"
	"fmt"
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

// NewKeystoreWallet returns new instance of KeystoreWallet
func NewKeystoreWallet(ks *keystore.KeyStore, accountHex, keystorePassphrase string) (ksw *KeystoreWallet, err error) {
	ksw = new(KeystoreWallet)
	a1 := accounts.Account{}
	a1.Address = common.HexToAddress(accountHex)
	account, err := ks.Find(a1)
	if err != nil {
		log.Fatalln("NewKeystoreWallet - Account not found within keystore", err.Error())
	}
	for _, w := range ks.Wallets() {
		for _, acc := range w.Accounts() {
			if acc.Address.Hash() == account.Address.Hash() {
				ksw.Wallet = w
				break
			}
		}
	}
	err = ks.Unlock(account, keystorePassphrase)
	if err != nil {
		log.Fatalln("NewKeystoreWallet - ", account.Address.Hex(), " could not be unlock: ", err.Error())
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

// GetNonceNumber gets actual nonce number of an Ethereum address/account
func (w *KeystoreWallet) GetNonceNumber(client *ethclient.Client) (nonce uint64, err error) {
	err = nil
	nonce, err = client.PendingNonceAt(context.Background(), w.Account.Address)
	if err != nil {
		return
	}
	return
}

func (w *KeystoreWallet) SwitchAccount(accountHex, keystorePassphrase string) (err error) {
	err = nil
	a1 := accounts.Account{}
	a1.Address = common.HexToAddress(accountHex)
	account, err := w.Keystore.Find(a1)
	if err != nil {
		err = fmt.Errorf("NewKeystoreWallet - SwitchAccount - Account: %s not found in keystore. Error: %s", accountHex, err.Error())
		return
	}
	var accountFound bool
	for _, wallet := range w.Keystore.Wallets() {
		for _, acc := range wallet.Accounts() {
			if acc.Address.Hash() == account.Address.Hash() {
				w.Wallet = wallet
				accountFound = true
				break
			}
		}
	}
	if !accountFound {
		err = fmt.Errorf("NewKeystoreWallet - SwitchAccount - Account: %s not found in wallets. Error: %s", accountHex, err.Error())
		return
	}
	w.Account = account

	err = w.Keystore.Unlock(account, keystorePassphrase)
	if err != nil {
		log.Fatalln("NewKeystoreWallet - ", account.Address.Hex(), " could not be unlock: ", err.Error())
	}
	return
}

/*
UpdateKeyedTransactor updates a keyed (signed?) transctor do perform a transaction within a Simulated Ethereum Blockchain
*/
func (w *KeystoreWallet) UpdateKeyedTransactor(transactor *bind.TransactOpts, client *ethclient.Client, increaseNonceFactor int, valueToSend int) (err error) {
	err = nil
	basicNonce, err := w.GetNonceNumber(client)
	if err != nil {
		log.Printf("[GetKeyedTransactor] Failure when get nonce from the network: %+v", err)
		return
	}
	nonce := basicNonce + uint64(increaseNonceFactor)

	latestEthBlockHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Println("[GetKeyedTransactorWithOptions] Error getting the latest Eth Block Header: ", err.Error())
		return
	}

	gasTip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Println("[GetKeyedTransactorWithOptions] Error getting SuggestGasTipCap: ", err.Error())
		return
	}

	maxGasFeeAccepted := new(big.Int).Add(
		latestEthBlockHeader.BaseFee,
		gasTip)

	transactor.GasLimit = uint64(6869310)
	transactor.Context = context.Background()
	transactor.GasFeeCap = maxGasFeeAccepted
	transactor.GasTipCap = gasTip
	transactor.Value = big.NewInt(int64(valueToSend))
	transactor.Nonce = big.NewInt(int64(nonce))
	return
}

// NewKeyStoreTransactor is a utility method to easily create a transaction signer from
// an decrypted key from a keystore
func (w *KeystoreWallet) NewKeyStoreTransactor(passphrase string, client *ethclient.Client) (*bind.TransactOpts, error) {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Println("[NewKeyStoreTransactor] Error getting chainID: ", err.Error())
		return nil, err
	}
	signer := types.NewLondonSigner(chainID)
	txOpts, err := bind.NewKeyStoreTransactorWithChainID(w.Keystore, w.Account, chainID)
	if err != nil {
		log.Println("[NewKeyStoreTransactor] Error generating NewKeyStoreTransactorWithChainID: ", err.Error())
		return nil, err
	}
	txOpts.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if address != w.Account.Address {
			return nil, errors.New("not authorized to sign this account")
		}
		signature, err := w.Keystore.SignHashWithPassphrase(w.Account, passphrase, signer.Hash(tx).Bytes())
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}

	return txOpts, nil
}

// GenerateSignedTxAsJSON generates a JSON signed raw Ethereum Transaction
// TODO: Update to London fork
func (w *KeystoreWallet) GenerateSignedTxAsJSON(
	txOpts *bind.TransactOpts,
	passphrase string,
	contractMethodParameters []byte,
	smartContractAddress common.Address,
	nonce, chainID uint64,
) (txJSON []byte, err error) {
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

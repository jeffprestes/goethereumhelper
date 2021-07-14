package goethereumhelper

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetMockBlockchain get a "in-memory" Blockchain instance
func GetMockBlockchain() (auth *bind.TransactOpts, backend *backends.SimulatedBackend, coinbaseAccountPrivateKey *ecdsa.PrivateKey) {
	coinbaseAccountPrivateKey, _ = crypto.GenerateKey()
	auth = bind.NewKeyedTransactor(coinbaseAccountPrivateKey)
	auth.Context = context.Background()
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(9000000000000000)}
	backend = backends.NewSimulatedBackend(alloc, 90000000)
	return
}

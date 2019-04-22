package goethereumhelper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

//GetMockBlockchain get a "in-memory" Blockchain instance
func GetMockBlockchain() (auth *bind.TransactOpts, backend *backends.SimulatedBackend) {
	key, _ := crypto.GenerateKey()
	auth = bind.NewKeyedTransactor(key)
	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(9000000000000000)}
	backend = backends.NewSimulatedBackend(alloc, 90000000)
	return
}

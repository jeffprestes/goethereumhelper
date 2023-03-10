package goethereumhelper

import (
	"context"
	"log"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// SubLogs Subscribe to watch to notifications to a specific Ethereum address
func SubLogs(addressToWatch common.Address, wg *sync.WaitGroup) {
	log.Println("[SubLogs] Waiting network confirmation to address ", addressToWatch.String(), " ...")
	defer wg.Done()
	wsClient, err := GetCustomNetworkClientWebsocket("<<put in here your EVM Node URL>>")
	if err != nil {
		log.Printf("[SubLogs] Houve falha ao conectar via WS na rede Rinkeby: %+v", err)
		return
	}
	query := ethereum.FilterQuery{
		Addresses: []common.Address{addressToWatch},
	}
	logs := make(chan types.Log)
	sub, err := wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("[SubLogs] Houve falha ao assinar eventos na rede Rinkeby para acompanhar o deploy do contrato na rede Rinkeby: %+v", err)
		return
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("[SubLogs] Erro recebido da rede Rinkeby: %+v", err)
			return
		case infoLog := <-logs:
			log.Println("[SubLogs] Information received from Rinkeby. Address: ", addressToWatch.String(), " - Information: ", infoLog)
		}
	}
}

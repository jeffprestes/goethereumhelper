package goethereumhelper

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

/*
GetRinkebyClient connects and return a client to the Rinkeby Ethereum network
*/
func GetRinkebyClient() (client *ethclient.Client, err error) {
	err = nil
	client, err = ethclient.Dial("https://rinkeby.infura.io/QPF0qjGpH9OjFuuMrCse")
	if err != nil {
		log.Printf("Houve falha ao conectar com Rinkeby via infuria: %+v", err)
		return
	}
	return
}

/*
GetRinkebyClientWebsocket connects via websocket and return a client to the Rinkeby Ethereum network
*/
func GetRinkebyClientWebsocket() (client *ethclient.Client, err error) {
	err = nil
	client, err = ethclient.Dial("wss://rinkeby.infura.io/ws")
	if err != nil {
		log.Printf("Houve falha ao conectar com Rinkeby via infuria: %+v", err)
		return
	}
	return
}

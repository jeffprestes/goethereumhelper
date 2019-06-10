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
	client, err = ethclient.Dial("http://rinkeby.caralabs.me:18575")
	if err != nil {
		log.Printf("Houve falha ao conectar com Rinkeby via Caralabs: %+v", err)
		return
	}
	return
}

/*
GetRinkebyClientWebsocket connects via websocket and return a client to the Rinkeby Ethereum network
*/
func GetRinkebyClientWebsocket() (client *ethclient.Client, err error) {
	err = nil
	client, err = ethclient.Dial("ws://rinkeby.caralabs.me:18576")
	if err != nil {
		log.Printf("Houve falha ao conectar com Rinkeby usando Websocket via Caralabs: %+v", err)
		return
	}
	return
}

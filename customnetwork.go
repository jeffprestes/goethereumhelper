package goethereumhelper

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

/*
GetCustomNetworkClient connects and return a client to user defined Ethereum network
*/
func GetCustomNetworkClient(URL string) (client *ethclient.Client, err error) {
	err = nil
	client, err = ethclient.Dial(URL)
	if err != nil {
		log.Printf("There was a failure connecting to %s: %+v", URL, err)
		return
	}
	return
}

/*
GetCustomNetworkClient connects via websocket and return a client to user defined Ethereum network
*/
func GetRinkebyClientWebsocket(URL string) (client *ethclient.Client, err error) {
	err = nil
	client, err = ethclient.Dial(URL)
	if err != nil {
		log.Printf("There was a failure connecting to %s via Websocket: %+v", URL, err)
		return
	}
	return
}

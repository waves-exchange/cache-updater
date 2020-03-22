package entities

import (
	"encoding/json"
	"fmt"
	"github.com/ventuary-lab/cache-updater/swagger-types/models"
	"io/ioutil"
	"net/http"
	"os"
)

func FetchLastTransactions (address string, lastCount uint, afterId *string) *[][]models.Transaction {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf(
		"%v/transactions/address/%v/limit/%v",
		nodeUrl, address, fmt.Sprintf("%v", lastCount),
	)

	if afterId != nil {
		connectionUrl += fmt.Sprintf("?after=%v", *afterId)
	}
	response, err := http.Get(connectionUrl)

	var txList [][]models.Transaction

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &txList
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return &txList
	}

	_ = json.Unmarshal(byteValue, &txList)
	return &txList
}

func FetchLastBlock () []byte {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/headers/last", nodeUrl)
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return make([]byte, 0)
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return make([]byte, 0)
	}

	return byteValue
}

func FetchBlocksRangeByAddress(address, heightMin, heightMax string) *[]models.Block {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/address/%v/%v/%v", nodeUrl, address, heightMin, heightMax)
	response, err := http.Get(connectionUrl)

	var blocksRange []models.Block

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &blocksRange
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return &blocksRange
	}

	_ = json.Unmarshal(byteValue, &blocksRange)
	return &blocksRange
}

func FetchTransactionsOnSpecificBlock (height string) *models.Block {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/at/%v", nodeUrl, height)
	response, err := http.Get(connectionUrl)

	var block models.Block

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &block
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return &block
	}

	_ = json.Unmarshal(byteValue, &block)

	return &block
}

func FetchBlocksRange (heightMin, heightMax string) *[]models.BlockHeader {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/headers/seq/%v/%v", nodeUrl, heightMin, heightMax)
	response, err := http.Get(connectionUrl)

	var blocksRange []models.BlockHeader

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &blocksRange
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return &blocksRange
	}

	_ = json.Unmarshal(byteValue, &blocksRange)

	return &blocksRange
}

func FetchStateChanges (txId string) *models.StateChanges {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/debug/stateChanges/info/%v", nodeUrl, txId)
	response, err := http.Get(connectionUrl)

	var stateChanges models.StateChanges

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &stateChanges
	}

	defer response.Body.Close()

	byteValue, readErr := ioutil.ReadAll(response.Body)

	if readErr != nil {
		return &stateChanges
	}
	_ = json.Unmarshal(byteValue, &stateChanges)

	return &stateChanges
}
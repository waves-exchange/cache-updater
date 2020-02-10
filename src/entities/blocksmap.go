package entities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const BLOCKS_MAP_NAME = "blocks_map"

type BlocksMap struct {
	tableName struct{} `pg:"blocks_map"`

	Height, Timestamp uint64
}

func (this *BlocksMap) GetBlocksMapByHeight (height string) *BlocksMap {
	// http://nodes.wavesplatform.com/blocks/at/77777
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := nodeUrl + "/blocks/at/" + height
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &BlocksMap{}
	}

	defer response.Body.Close()
	
	byteValue, _ := ioutil.ReadAll(response.Body)

	var blocksMap BlocksMap
	json.Unmarshal([]byte(byteValue), &blocksMap)

	return &blocksMap
}

func (this *BlocksMap) GetBlocksMapSequenceByRange (heightMin, heightMax string) *[]BlocksMap {
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := fmt.Sprintf("%v/blocks/headers/seq/%v/%v", nodeUrl, heightMin, heightMax)
	response, err := http.Get(connectionUrl)

	var blocksMapArray []BlocksMap

	if err != nil {
		fmt.Printf("Error occured on fetch... %v \n", err)
		return &blocksMapArray
	}

	defer response.Body.Close()

	byteValue, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(byteValue), &blocksMapArray)

	return &blocksMapArray
}
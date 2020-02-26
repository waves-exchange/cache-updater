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

func (bm *BlocksMap) GetBlocksMapByHeight (height string) *BlocksMap {
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
	_ = json.Unmarshal([]byte(byteValue), &blocksMap)

	return &blocksMap
}

func (bm *BlocksMap) GetBlocksMapSequenceByRange (heightMin, heightMax string) *[]BlocksMap {
	var blocksMapArray []BlocksMap

	byteValue := FetchBlocksRange(heightMin, heightMax)
	_ = json.Unmarshal([]byte(byteValue), &blocksMapArray)

	return &blocksMapArray
}
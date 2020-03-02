package controllers

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/entities"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type UpdateController struct {
	DbDelegate *DbController
	ScDelegate *ShareController
}

func (uc *UpdateController) GrabAllAddressData () ([]byte, error) {
	dAppAddress := os.Getenv("DAPP_ADDRESS")
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := nodeUrl + "/addresses/data/" + dAppAddress
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Println(err)
		return make([]byte, 0), err
	}

	defer response.Body.Close()

	byteValue, _ := ioutil.ReadAll(response.Body)

	return byteValue, nil
}

func (uc *UpdateController) UpdateStateChangedData (latestExRecord entities.BondsOrder, maxHeightRange uint64) {
	minH := latestExRecord.Height
	maxH := minH + maxHeightRange

	blocks := entities.FetchBlocksRange(
		fmt.Sprintf("%v", minH),
		fmt.Sprintf("%v", maxH),
	)

	//mutableKeys := []string { entities.OrderStatusKey, entities.OrderFilledTotalKey }
	//staticKeys := []string { entities.OrderBookKey, entities.OrderFirstKey }
	//joinedMutableKeys := strings.Join(mutableKeys, "_")
	delimiter := "_"

	for _, block := range *blocks {
		blockWithTxList := entities.FetchTransactionsOnSpecificBlock(
			fmt.Sprintf("%v", *block.Height),
		)

		// Invoke Script Transaction: 16
		for _, tx := range blockWithTxList.Transactions {
			txType := tx["type"]
			//fmt.Printf("Type is: %v \n", txType)

			// Let only Invoke transactions stay
			if txType != float64(16) {
				continue
			}

			txId := tx["id"]
			txSender := tx["sender"].(string)

			wrappedStateChanges := entities.FetchStateChanges(txId.(string))

			stateChanges := wrappedStateChanges.StateChanges

			if !(stateChanges.Data != nil && len(stateChanges.Data) > 0) {
				return
			}

			fmt.Printf("TxId: %v Len: %v; StateChange: %+v \n", txId, len(stateChanges.Data), *stateChanges.Data[0])

			if txSender == "" {
				continue
			}
			if len(stateChanges.Data) != 12 {
				continue
			}

			for i, change := range stateChanges.Data {
				changeKey := *(*change).Key

				if changeKey == entities.OrderBookKey || changeKey == entities.OrderFirstKey {
					continue
				}

				if !strings.Contains(changeKey, "order") {
					break
				}

				splittedKey := strings.Split(changeKey, delimiter)
				if len(splittedKey) < 3 {
					continue
				}

				orderId := splittedKey[len(splittedKey) - 1]
				dict := entities.MapStateChangesDataToDict(stateChanges)
				fmt.Printf("Dict: %v \n", dict)
				entity := uc.ScDelegate.BondsOrder.MapItemToModel(orderId, dict)

				//updateErr := uc.DbDelegate.DbConnection.Update(entity)
				//
				//if updateErr != nil {
				//	fmt.Printf("Update result: %v \n", updateErr)
				//}

				fmt.Printf("Entity: %+v \n", entity)
				fmt.Printf("Data key immutable part: %v \n", changeKey)
				fmt.Printf("TX ID: %v, Sender is: %v \n", txId, txSender)

				fmt.Printf("Data #%v: %+v \n", i + 1, *change)
				fmt.Printf("%v , %v , %v \n", *(*change).Key, (*change).Value, *(*change).Type)

				break
			}
		}
	}
}

func (uc *UpdateController) UpdateAllData () {
	uc.DbDelegate.HandleRecordsUpdate()
}
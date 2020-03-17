package controllers

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/entities"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func (uc *UpdateController) UpdateStateChangedData (
	minHeight, maxHeight uint64,
) {
	blocks := entities.FetchBlocksRange(
		fmt.Sprintf("%v", minHeight),
		fmt.Sprintf("%v", maxHeight),
	)

	delimiter := "_"

	for _, block := range *blocks {
		blockWithTxList := entities.FetchTransactionsOnSpecificBlock(
			fmt.Sprintf("%v", *block.Height),
		)

		fmt.Printf("Checking block: %v \n", *block.Height)

		// Invoke Script Transaction: 16
		for _, tx := range blockWithTxList.Transactions {
			txType := tx["type"]

			// Let only Invoke transactions stay
			if txType != float64(16) {
				continue
			}

			txId := tx["id"]
			txSender := tx["sender"].(string)
			txDapp := tx["dApp"].(string)
			txCall := tx["call"].(map[string]interface{})
			txCallFunction := txCall["function"]

			if txDapp != os.Getenv("DAPP_ADDRESS") {
				continue
			}

			mutateMethodNames := []string{ "addBuyBondOrder", "sellBond", "cancelOrder" }

			isMutateMethod := false
			for _, methodName := range mutateMethodNames {
				if methodName == txCallFunction {
					isMutateMethod = true
					break
				}
			}
			if !isMutateMethod { continue }

			wrappedStateChanges := entities.FetchStateChanges(txId.(string))

			stateChanges := wrappedStateChanges.StateChanges

			if !(stateChanges.Data != nil && len(stateChanges.Data) > 0) {
				return
			}

			fmt.Printf("TX: %v L: %v; StateChange: %v \n", txId, len(stateChanges.Data), *stateChanges.Data[0])

			for _, change := range stateChanges.Data {
				changeKey := *(*change).Key

				if *block.Height == 1974626 {
					fmt.Printf("CHANGE KEY ON BLOCK %v is %v",1974626, changeKey)
				}

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
				entity := uc.ScDelegate.BondsOrder.MapItemToModel(orderId, dict)
				fmt.Printf("Entity: %+v \n", entity)
				fmt.Printf("Order ID IS: %v \n", entity.OrderId)

				exists, _ := uc.DbDelegate.DbConnection.Model(entity).Where("order_id = ?", entity.OrderId).Exists()

				if exists {
					_, updateErr := uc.DbDelegate.DbConnection.Model(entity).
						Where("order_id = ?", entity.OrderId).
						Update(entity)

					entities.DefaultErrorHandler(updateErr)
				} else {
					insertErr := uc.DbDelegate.DbConnection.Insert(entity)
					entities.DefaultErrorHandler(insertErr)
				}

				break
			}
		}
	}
}

func (uc *UpdateController) UpdateAllData () {
	uc.DbDelegate.HandleRecordsUpdate()
}

func (uc *UpdateController) StartConstantUpdating () {
	frequency := os.Getenv("UPDATE_FREQUENCY")
	duration, convErr := strconv.Atoi(frequency)

	entities.DefaultErrorHandler(convErr)
	if convErr != nil {
		return
	}
	duration = int(time.Duration(duration) * time.Millisecond)

	for {
		uc.UpdateAllData()
		time.Sleep(time.Duration(duration))
	}
}
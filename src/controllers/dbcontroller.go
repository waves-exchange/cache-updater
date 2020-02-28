package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/ventuary-lab/cache-updater/src/entities"
	//"reflect"
	//"os"
	//"reflect"
	"strconv"
)

type DbController struct {
	UcDelegate *UpdateController
	DbConnection *pg.DB
}

func (dc *DbController) ConnectToDb () {
	dbhost, dbport, dbuser, dbpass, dbdatabase := entities.GetDBCredentials()

	db := pg.Connect(&pg.Options{
		Addr: dbhost + ":" + dbport,
		User:     dbuser,
		Password: dbpass,
		Database: dbdatabase,
	})

	dc.DbConnection = db
}

func (dc *DbController) HandleRecordsUpdate () {
	var records []entities.DAppStringRecord
	var numberRecords []entities.DAappNumberRecord

	var existingBondsOrders []entities.BondsOrder
	_ = dc.GetAllEntityRecords(&existingBondsOrders, entities.BONDS_ORDERS_NAME)

	if len(existingBondsOrders) != 0 {
		dc.HandleExistingBondsOrdersUpdate()
		return
	}

	byteValue, _ := dc.UcDelegate.GrabAllAddressData()

	json.Unmarshal([]byte(byteValue), &records)
	json.Unmarshal([]byte(byteValue), &numberRecords)

	nodeData := map[string]string{}

	for i := 0; i < len(records); i++ {
		record := records[i]

		if *record.Value == "" {
			numberRecord := numberRecords[i]

			*record.Value = strconv.Itoa(*numberRecord.Value)
		}

		nodeData[record.Key] = *record.Value
	}

	var orderheights []uint64
	bondsorders := dc.UcDelegate.ScDelegate.BondsOrder.UpdateAll(&nodeData)

	dc.HandleBondsOrdersUpdate(&bondsorders)

	for _, order := range bondsorders {
		orderheights = append(orderheights, order.Height)
	}

	dc.HandleBlocksMapUpdate(&orderheights)
}

func (dc *DbController) GetAllEntityRecords (records interface{}, tableName string) error {
	_, getRecordsErr := dc.DbConnection.
		Query(records, fmt.Sprintf("SELECT * FROM %v;", tableName))
	return getRecordsErr
}

func (dc *DbController) HandleExistingBondsOrdersUpdate () {
	fmt.Println("Records exists, updating based on existing...")

	var records []entities.BondsOrder
	_, getRecordsErr := dc.DbConnection.
		Query(&records, fmt.Sprintf("SELECT * FROM %v ORDER BY height DESC;", entities.BONDS_ORDERS_NAME))

	if getRecordsErr != nil {
		fmt.Printf("Error on select... %v\n", getRecordsErr)
	}

	var bm entities.BlocksMap
	latestExRecord := records[0]
	byteValue := entities.FetchLastBlock()
	_ = json.Unmarshal([]byte(byteValue), &bm)

	maxHeightRange := uint64(99)
	heightDiff := bm.Height - latestExRecord.Height

	fmt.Printf("heightDiff: %v\n", heightDiff)

	dc.UcDelegate.UpdateStateChangedData(latestExRecord, maxHeightRange)

	//if true {
	//	minH := latestExRecord.Height
	//	maxH := minH + maxHeightRange
	//
	//	blocks := entities.FetchBlocksRange(
	//		fmt.Sprintf("%v", minH),
	//		fmt.Sprintf("%v", maxH),
	//	)
	//
	//	mutableKeys := []string { entities.OrderStatusKey, entities.OrderFilledTotalKey }
	//	joinedMutableKeys := strings.Join(mutableKeys, "_")
	//	delimiter := "_"
	//
	//	for _, block := range *blocks {
	//		blockWithTxList := entities.FetchTransactionsOnSpecificBlock(
	//			fmt.Sprintf("%v", *block.Height),
	//		)
	//
	//		// Invoke Script Transaction: 16
	//		for _, tx := range blockWithTxList.Transactions {
	//			txType := tx["type"]
	//			//fmt.Printf("Type is: %v \n", txType)
	//
	//			// Let only Invoke transactions stay
	//			if txType != float64(16) {
	//				continue
	//			}
	//
	//			txId := tx["id"]
	//			txSender := tx["sender"].(string)
	//
	//			wrappedStateChanges := entities.FetchStateChanges(txId.(string))
	//
	//			stateChanges := wrappedStateChanges.StateChanges
	//
	//			if !(stateChanges.Data != nil && len(stateChanges.Data) > 0) {
	//				return
	//			}
	//
	//			//fmt.Printf("TxId: %v Len: %v; StateChange: %+v \n", txId, len(stateChanges.Data), *stateChanges.Data[0])
	//
	//			if txSender == "" {
	//				continue
	//			}
	//
	//			//fmt.Printf("Data: %+v \n", *stateChanges.Data[0])
	//
	//			for i, change := range stateChanges.Data {
	//				changeKey := *(*change).Key
	//
	//				splittedKey := strings.Split(changeKey, delimiter)
	//
	//				if len(splittedKey) < 3 {
	//					continue
	//				}
	//
	//				changeKey = strings.Join(splittedKey[:(len(splittedKey) - 1)], delimiter)
	//
	//				if !strings.Contains(joinedMutableKeys, changeKey) {
	//					continue
	//				}
	//
	//				fmt.Printf("Data key immutable part: %v \n", changeKey)
	//				fmt.Printf("Sender is: %v \n", txSender)
	//				fmt.Printf("Data #%v: %+v \n", i + 1, *change)
	//				fmt.Printf("%v , %v , %v \n", *(*change).Key, (*change).Value, *(*change).Type)
	//			}
	//		}
	//	}
	//
	//} else {
	//
	//}
}

func (dc *DbController) HandleBondsOrdersUpdate (freshData *[]*entities.BondsOrder) {
	var existingRecords []entities.BondsOrder
	getRecordsErr := dc.GetAllEntityRecords(&existingRecords, entities.BONDS_ORDERS_NAME)

	if getRecordsErr != nil {
		return
	}

	isEmpty := len(existingRecords) == 0

	// Base case when table is empty, just upload and return
	if !isEmpty {
		return
	}

	fmt.Printf("0 records exist \n")
	if len(*freshData) == 0 {
		fmt.Printf("0 new records added \n")
		return
	}

	insertErr := dc.DbConnection.Insert(freshData)

	if insertErr != nil {
		fmt.Printf("Error occured on Insert... %v \n", insertErr)
	} else {
		fmt.Printf("Successfully inserted %v rows \n", len(*freshData))
	}
}

func (dc *DbController) HandleBlocksMapUpdate (heightarr *[]uint64) {
	var existingRecords []entities.BlocksMap
	var bondsOrders []entities.BondsOrder

	_, getRecordsErr := dc.DbConnection.
		Query(&existingRecords, fmt.Sprintf("SELECT * FROM %v ORDER BY height ASC;", entities.BLOCKS_MAP_NAME))
	_, getBondsOrdersErr := dc.DbConnection.
		Query(&bondsOrders, fmt.Sprintf("SELECT height FROM %v GROUP BY height ORDER BY height ASC;", entities.BONDS_ORDERS_NAME))

	if getRecordsErr != nil || getBondsOrdersErr != nil {
		fmt.Printf("Error occured on Query Select... %v; %v \n", getRecordsErr, getBondsOrdersErr)
		return
	}

	if len(bondsOrders) == 0 {
		fmt.Printf("%v table is empty... \n", entities.BONDS_ORDERS_NAME)
		return
	}

	var freshData []entities.BlocksMap

	minHeightBm := bondsOrders[0]
	maxHeightBm := bondsOrders[len(bondsOrders) - 1]
	maxRecordsCount := uint64(99)
	bm := dc.UcDelegate.ScDelegate.BlocksMap

	if len(existingRecords) > 0 {
		minExRecord := existingRecords[len(existingRecords) - 1]
		minHeightBm = entities.BondsOrder{ Height: minExRecord.Height + 1, Timestamp: minExRecord.Timestamp }
	}

	minHeight := minHeightBm.Height
	maxHeight := minHeightBm.Height + maxRecordsCount
	index := 1
	iterationsLimitPerUpdate := 15

	for {
		fmt.Printf("min: %v, max: %v \n", minHeight, maxHeight)
		fetchedBlocksMap := bm.GetBlocksMapSequenceByRange(fmt.Sprintf("%v", minHeight), fmt.Sprintf("%v", maxHeight))

		freshData = append(freshData, *fetchedBlocksMap...)
		minHeight = maxHeight + 1
		maxHeight = maxHeight + maxRecordsCount + 1

		if maxHeight == maxHeightBm.Height {
			break
		}
		if maxHeight > maxHeightBm.Height {
			maxHeight = maxHeightBm.Height
		}

		index++

		if index == iterationsLimitPerUpdate {
			break
		}
	}

	fmt.Printf("Data len is: %v \n", len(freshData))

	fmt.Printf("blocks count: %v \n", len(freshData))
	insertErr := dc.DbConnection.Insert(&freshData)

	if insertErr != nil {
		fmt.Printf("Error occured on Insert... %v \n", insertErr)
	} else {
		fmt.Printf("Successfully inserted %v rows \n", len(freshData))
	}
}
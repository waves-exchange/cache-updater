package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/ventuary-lab/cache-updater/models"
	"github.com/ventuary-lab/cache-updater/src/entities"
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
		PoolSize: 1,
	})

	dc.DbConnection = db
}

func (dc *DbController) HandleRecordsUpdate () {
	byteValue, _ := dc.UcDelegate.GrabAllAddressData()
	nodeData := entities.MapNodeDataToDict(byteValue)

	bondsOrders := dc.UcDelegate.ScDelegate.BondsOrder.UpdateAll(&nodeData)
	neutrinoOrders := dc.UcDelegate.ScDelegate.NeutrinoOrder.UpdateAll(&nodeData)

	dc.HandleBondsOrdersUpdate(&bondsOrders)
	dc.HandleNeutrinoOrdersUpdate(&neutrinoOrders)
	dc.HandleBlocksMapUpdate()
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

	var lastBlockHeader models.BlockHeader
	latestExRecord := records[0]
	byteValue := entities.FetchLastBlock()
	_ = json.Unmarshal([]byte(byteValue), &lastBlockHeader)
	lastBlockHeaderHeight := uint64(*lastBlockHeader.Height)

	maxHeightRange := uint64(99)
	heightDiff := lastBlockHeaderHeight - latestExRecord.Height

	fmt.Printf("heightDiff: %v\n", heightDiff)
	minH := latestExRecord.Height
	maxH := minH + maxHeightRange

	if heightDiff > maxHeightRange {
		for {
			if minH > lastBlockHeaderHeight {
				break
			}

			dc.UcDelegate.UpdateStateChangedData(minH, maxH)
			minH = maxH + 1
			maxH = minH + maxHeightRange
		}
	} else {
		dc.UcDelegate.UpdateStateChangedData(minH, maxH)
	}
}

func (dc *DbController) HandleNeutrinoOrdersUpdate (freshData *[]*entities.NeutrinoOrder) {
	if len(*freshData) == 0 {
		fmt.Printf("0 new records added \n")
		return
	}

	_, dropRecordsErr := dc.DbConnection.
		Query(nil, fmt.Sprintf("DELETE FROM %v;", entities.NEUTRINO_ORDERS_NAME))
	if dropRecordsErr != nil {
		fmt.Printf("Error on drop all neutrino orders records... %v;", dropRecordsErr)
		return
	}
	
	insertErr := dc.DbConnection.Insert(freshData)

	if insertErr != nil {
		fmt.Printf("Error occured on Insert... %v \n", insertErr)
	} else {
		fmt.Printf("Successfully inserted %v rows \n", len(*freshData))
	}
}

func (dc *DbController) HandleBondsOrdersUpdate (freshData *[]*entities.BondsOrder) {
	if len(*freshData) == 0 {
		fmt.Printf("0 new records added \n")
		return
	}

	_, dropRecordsErr := dc.DbConnection.
		Query(nil, fmt.Sprintf("DELETE FROM %v;", entities.BONDS_ORDERS_NAME))
	if dropRecordsErr != nil {
		fmt.Printf("Error on drop all bonds orders records... %v;", dropRecordsErr)
		return
	}

	insertErr := dc.DbConnection.Insert(freshData)

	if insertErr != nil {
		fmt.Printf("Error occured on Insert... %v \n", insertErr)
	} else {
		fmt.Printf("Successfully inserted %v rows \n", len(*freshData))
	}
}

func (dc *DbController) GetEntityRecordsCount (entity interface{}) (error, int) {
	var count int
	_, err := dc.DbConnection.Model(entity).QueryOne(pg.Scan(&count), `
	    SELECT count(*)
	    FROM ?TableName AS ?TableAlias
	`)

	return err, count
}

func (dc *DbController) HandleBlocksMapUpdate () {
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

		if minHeight > maxHeight {
			break
		}

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

	fmt.Printf("blocks count: %v \n", len(freshData))
	insertErr := dc.DbConnection.Insert(&freshData)

	t := time.Now().Format(time.RFC3339)

	if insertErr != nil {
		fmt.Printf("%v; Error occured on Insert... %v \n", t, insertErr)
	} else {
		fmt.Printf("%v; bonds_orders: Successfully inserted %v rows \n", t, len(freshData))
	}
}
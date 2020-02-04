package controllers;

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/entities"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"os"
	"encoding/json"
	"strconv"
	// "regexp"
)

type DbController struct {
	UcDelegate *UpdateController
	DbConnection *pg.DB
}

func (dbc *DbController) ConnectToDb () {
	envLoadErr := godotenv.Load(".env")
	if envLoadErr != nil {
		_ = godotenv.Load(".env.example")
	}

	dbuser := os.Getenv("DB_USERNAME");
	dbpass := os.Getenv("DB_PASS");
	dbdatabase := os.Getenv("DB_NAME");
  
	db := pg.Connect(&pg.Options{
		User:     dbuser,
		Password: dbpass,
		Database: dbdatabase,
	})

	dbc.DbConnection = db
}

func (db *DbController) HandleRecordsUpdate (byteValue []byte) {
	var records []entities.DAppStringRecord;
	var numberRecords []entities.DAappNumberRecord;

	json.Unmarshal([]byte(byteValue), &records)
	json.Unmarshal([]byte(byteValue), &numberRecords)

	nodeData := map[string]string{};

	for i := 0; i < len(records); i++ {
		record := records[i]

		if *record.Value == "" {
			numberRecord := numberRecords[i]

			*record.Value = strconv.Itoa(*numberRecord.Value)
		}

		nodeData[record.Key] = *record.Value;
	}

	var bondsorders []entities.BondsOrder

	rawbo := entities.BondsOrder{}
	bondsorders = rawbo.UpdateAll(&nodeData)

	db.HandleBondsOrdersUpdate(&bondsorders)
}

func (db *DbController) HandleBondsOrdersUpdate (freshData *[]entities.BondsOrder) {
	// fmt.Printf("ORDER: \n %+v \n", freshData[0])
	var existingRecords []entities.BondsOrder

	_, getRecordsErr := db.DbConnection.Query(&existingRecords, "SELECT * FROM f_bonds_orders;")

	// fmt.Printf("getRecordsErr... %v \n records len: %v \n", getRecordsErr, len(existingRecords))
	if getRecordsErr != nil {
		return
	}

	// Base case when table is empty, just upload and return
	if len(existingRecords) == 0 {
		fmt.Printf("0 records exist \n")
		insertErr := db.DbConnection.Insert(freshData)

		if insertErr != nil {
			fmt.Printf("Error occured on Insert... %v \n", insertErr)
		} else {
			fmt.Printf("Successfully inserted %v rows \n", len(*freshData))
		}

		return
	}

	// Case with diff between two data sets

	var recordsToUpdate []entities.BondsOrder;
	// var recordsToInsert []entities.BondsOrder;

	// isSameLength := len(existingRecords) == len(*freshData)
	// isNewLonger := len(*freshData) > len(existingRecords)

	for _, newRecord := range *freshData {
		for _, existingRecord := range existingRecords {
			if newRecord != existingRecord {
				recordsToUpdate = append(recordsToUpdate, newRecord)
				// if isSameLength {
				// 	continue;
				// }
			}
		}
	}

	insertErr := db.DbConnection.Update(&recordsToUpdate)

	if insertErr != nil {
		fmt.Printf("InsertErr: %v \n", insertErr)
	}
}

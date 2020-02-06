package controllers;

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"github.com/ventuary-lab/cache-updater/src/entities"
)

type DbController struct {
	UcDelegate *UpdateController
	DbConnection *pg.DB
}

func (this *DbController) ConnectToDb () {
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

	this.DbConnection = db
}

func (this *DbController) HandleRecordsUpdate (byteValue []byte) {
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

	this.HandleBondsOrdersUpdate(&bondsorders)
}

func (this *DbController) HandleBondsOrdersUpdate (freshData *[]entities.BondsOrder) {
	var existingRecords []entities.BondsOrder

	_, getRecordsErr := this.DbConnection.Query(&existingRecords, "SELECT * FROM f_bonds_orders;")

	// fmt.Printf("getRecordsErr... %v \n records len: %v \n", getRecordsErr, len(existingRecords))
	if getRecordsErr != nil {
		return
	}

	isEmpty := len(existingRecords) == 0
	isNewLonger := len(*freshData) > len(existingRecords)

	// Base case when table is empty, just upload and return
	if isEmpty {
		fmt.Printf("0 records exist \n")
		insertErr := this.DbConnection.Insert(freshData)

		if insertErr != nil {
			fmt.Printf("Error occured on Insert... %v \n", insertErr)
		} else {
			fmt.Printf("Successfully inserted %v rows \n", len(*freshData))
		}
	} else {
		if isNewLonger {
			recordsToAdd := []entities.BondsOrder{}
			recordsToUpdate := []entities.BondsOrder{}
			bo := entities.BondsOrder{}

			for _, newRecord := range *freshData {
				if !bo.Includes(&existingRecords, &newRecord) {
					recordsToAdd = append(recordsToAdd, newRecord)
				} else {
					recordsToUpdate = append(recordsToAdd, newRecord)
				}
			}

			this.DbConnection.Update(&recordsToUpdate)
			_ = this.DbConnection.Insert(&recordsToAdd)

			fmt.Printf("Added %v rows; updated %v rows \n", len(recordsToAdd), len(recordsToUpdate))
		} else {
			this.DbConnection.Update(freshData)

			fmt.Printf("Successfully updated %v rows \n", len(*freshData))
		}
	}
}

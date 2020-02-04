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
	// _, err := db.DbConnection.Query(&bondsorders, `SELECT * FROM f_bonds_orders`)

	// if err == nil {
	// 	fmt.Printf("Len: %v \n", len(bondsorders))
	// 	// fmt.Printf("ORDER: \n %+v \n", bondsorders[0])
	// }

	rawbo := entities.BondsOrder{}
	bondsorders = rawbo.UpdateAll(&nodeData)

	fmt.Printf("ORDER: \n %+v \n", bondsorders[0])

	insertErr := db.DbConnection.Insert(&bondsorders[0])

	if insertErr != nil {
		fmt.Println("error occured", insertErr)
	}
}
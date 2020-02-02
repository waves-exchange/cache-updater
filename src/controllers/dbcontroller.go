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

	defer db.Close()

	var bondsorders []entities.NeutrinoOrder
	_, err := db.Query(&bondsorders, `SELECT * FROM neutrino_orders`)

	fmt.Println(bondsorders[0].GetKeys("1"))
	
	rawItem := map[string]string {
		"order_height_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ": "1868066",
		"order_owner_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ": "3PGmja5rWBPiQ7n9eLSBgQBd6EzTmUFgddB",
		"order_price_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ": "55",
		"order_total_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ": "1434000000",
		"order_status_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ": "canceled",
		"orderbook": "",
	};

	rawbo := entities.BondsOrder{}

	// fmt.Println("RAW ITEM: ", rawItem)
	fmt.Printf("RAWBO VALS: %v \n", rawbo.MapItemToModel("zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ", rawItem))

	if err != nil {
		fmt.Println("Error: ", err)
		return;
	}

	fmt.Println(bondsorders[len(bondsorders)-1])
}

func (db *DbController) HandleRecordsUpdate (byteValue []byte) {
	var records []entities.DAppStringRecord;
	var numberRecords []entities.DAappNumberRecord;

	json.Unmarshal([]byte(byteValue), &records)
	json.Unmarshal([]byte(byteValue), &numberRecords)

	maxcount := 10
	nodeData := map[string]string{};

	for i := 0; i < len(records); i++ {
		if i == maxcount {
			break
		}

		record := records[i]

		if *record.Value == "" {
			numberRecord := numberRecords[i]

			*record.Value = strconv.Itoa(*numberRecord.Value)
		}

		nodeData[record.Key] = *record.Value;
	}

	// fmt.Printf("Iteration ended. Example val %v \n", nodeData)

	bo := entities.BondsOrder{}
	bo.UpdateAll(&nodeData)
}
package controllers;

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/entities"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"os"
)

type DbController struct {}

func (dbc DbController) ConnectToDb () {
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

	// newRecord := entities.NeutrinoOrder {
	// 	Height: "135", Currency: "usd-nb", Owner: "sdfg", Total: "11", Ordernext: nil, Orderprev: nil, Order_id: "dfgdg",
	// 	Timestamp: 3563456,
	// 	Status: enums.NEW,
	// 	Resttotal: 154,
	// 	Type: enums.LIQUIDATE,
	// 	Isfirst: false, Islast: false,
	// }
	// insertErr := db.Insert(&newRecord)

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
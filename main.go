package main;

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/structs"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"os"
)

func connectToDb () {
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

	// var bondsorders []structs.BondsOrder
	// _, err := db.Query(&bondsorders, `SELECT * FROM bonds_orders`)

	newRecord := structs.NeutrinoOrder {
		Height: "135", Currency: "usd-nb", Owner: "sdfg", Total: "11", Ordernext: nil, Orderprev: nil, Order_id: "dfgdg",
		Timestamp: 3563456,
		Status: "new",
		Resttotal: 154,
		Type: "liquidate",
		Isfirst: false, Islast: false,
	}

	insertErr := db.Insert(&newRecord)

	var bondsorders []structs.NeutrinoOrder
	_, err := db.Query(&bondsorders, `SELECT * FROM neutrino_orders`)

	if insertErr != nil {
		fmt.Println("insertErr: ", insertErr)
		return;
	}
	
	if err != nil {
		fmt.Println("Error: ", err)
		return;
	}

	fmt.Println(bondsorders[len(bondsorders)-1])
}

func main () {
	connectToDb()
}
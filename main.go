package main;

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/entities"
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

	// if insertErr != nil {
	// 	fmt.Println("insertErr: ", insertErr)
	// 	return;
	// }
	
	if err != nil {
		fmt.Println("Error: ", err)
		return;
	}

	fmt.Println(bondsorders[len(bondsorders)-1])
}

func main () {
	connectToDb()
}
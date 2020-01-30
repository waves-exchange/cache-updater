package main;

import (
	"fmt"
	"github.com/ventuary-lab/cache-updater/structs"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"os"
	// "github.com/go-pg/pg/v9/orm"
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

	var bondsorders []structs.BondsOrder
	_, err := db.Query(&bondsorders, `SELECT * FROM bonds_orders`)
	
	if err != nil {
		fmt.Println("Error: ", err)
		return;
	}

	fmt.Println(bondsorders[0])
}

func main () {
	connectToDb()
}
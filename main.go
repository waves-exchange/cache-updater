package main;

import (
	"github.com/ventuary-lab/cache-updater/src/controllers"
	// "github.com/go-pg/pg/v9"
)

func main () {
	dbc := controllers.DbController{}
	uc := controllers.UpdateController{}
	dbc.UcDelegate = &uc
	uc.DbDelegate = &dbc

	dbc.ConnectToDb()
	uc.UpdateAllData()
}
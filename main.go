package main

import (
	"github.com/ventuary-lab/cache-updater/src/controllers"
	"github.com/ventuary-lab/cache-updater/src/entities"
)

func main () {
	dbc := &controllers.DbController{}
	uc := &controllers.UpdateController{}

	shc := &controllers.ShareController{}
	shc.BondsOrder = &entities.BondsOrder{}
	shc.BlocksMap = &entities.BlocksMap{}

	dbc.UcDelegate = uc
	uc.DbDelegate = dbc
	uc.ScDelegate = shc

	dbc.ConnectToDb()

	//uc.UpdateAllData()
	uc.StartConstantUpdating()
}
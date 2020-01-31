package main;

import (
	"github.com/ventuary-lab/cache-updater/src/controllers"
)

func main () {
	uc := controllers.DbController{}

	uc.ConnectToDb()
}
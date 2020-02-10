package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type UpdateController struct {
	DbDelegate *DbController
}

func (uc *UpdateController) UpdateAllData () {
	dAppAddress := os.Getenv("DAPP_ADDRESS")
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := nodeUrl + "/addresses/data/" + dAppAddress
	response, err := http.Get(connectionUrl)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	
	defer response.Body.Close()

	byteValue, _ := ioutil.ReadAll(response.Body)

	uc.DbDelegate.HandleRecordsUpdate(byteValue)
}
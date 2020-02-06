package controllers;

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
)

type UpdateController struct {
	DbDelegate *DbController
}

func (this *UpdateController) UpdateAllData () {
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

	this.DbDelegate.HandleRecordsUpdate(byteValue)
}
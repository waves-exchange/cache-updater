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

func (uc *UpdateController) GrabAllAddressData () ([]byte, error) {
	dAppAddress := os.Getenv("DAPP_ADDRESS")
	nodeUrl := os.Getenv("NODE_URL")
	connectionUrl := nodeUrl + "/addresses/data/" + dAppAddress
	response, err := http.Get(connectionUrl)

	if err != nil {
		fmt.Println(err)
		return make([]byte, 0), err
	}

	defer response.Body.Close()

	byteValue, _ := ioutil.ReadAll(response.Body)

	return byteValue, nil
}

// func (uc *UpdateController) GrabStateChangeData () ([]byte, error) {}

func (uc *UpdateController) UpdateAllData () {
	uc.DbDelegate.HandleRecordsUpdate()
}
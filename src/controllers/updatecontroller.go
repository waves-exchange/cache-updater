package controllers;

import (
	"fmt"
	"io/ioutil"
	"os"
)

type UpdateController struct {
	DbDelegate *DbController
}

func (uc *UpdateController) UpdateAllData () {
	jsonFile, err := os.Open("scriptdata.test.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	uc.DbDelegate.HandleRecordsUpdate(byteValue)
}
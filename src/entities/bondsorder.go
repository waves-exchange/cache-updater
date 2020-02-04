package entities;

import (
	"strconv"
	// "math"
	"regexp"
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/constants"
)

type BondsOrder struct {
	Height, Owner, Status, Index, Pairname, Type, Uuid string
	Price int
	Timestamp int64
	Total, Filledamount, Filledtotal, Resttotal, Amount, Restamount float64
}

func (bo *BondsOrder) GetKeys(id string) []string {
	return []string {
		"order_height_" + id,
		"order_owner_" + id,
		"order_price_" + id,
		"order_total_" + id,
		"order_filled_total_" + id,
		"order_status_" + id,
		"orderbook",
	}
}

func (bo *BondsOrder) UpdateAll (nodeData *map[string]string) []BondsOrder {
	ids := []string{}
	result := []BondsOrder{}
	defaultRawRegex := "([A-Za-z0-9]{40,50})"
	regexKeys := bo.GetKeys(defaultRawRegex)
	heightKey := regexKeys[0]
	heightRegex, heightRegexErr := regexp.Compile(heightKey)
	nodeKeys := []string{}
	resolveData := make(map[string](map[string]string))

	for k, _ := range *nodeData {
		for _, regexKey := range regexKeys {
			compiledRegex := regexp.MustCompile(regexKey)

			if len(compiledRegex.FindSubmatch([]byte(k))) == 0 {
				continue;
			}
		}

		nodeKeys = append(nodeKeys, k)
	}

	for _, k := range nodeKeys {
		heightRegexSubmatches := heightRegex.FindSubmatch([]byte(k))

		if len(heightRegexSubmatches) < 2 {
			continue
		}

		matchedAddress := string(heightRegexSubmatches[1])

		if matchedAddress != "" {
			ids = append(ids, matchedAddress)

			resolveData[matchedAddress] = map[string]string{}

			validKeys := bo.GetKeys(matchedAddress)

			for _, validKey := range validKeys {
				for _, k := range nodeKeys {
					if k == validKey {
						resolveData[matchedAddress][k] = (*nodeData)[k]
					}
				}
			}
		}
	}

	if heightRegexErr != nil {
		return result
	}

	raw := BondsOrder{}

	for _, id := range ids {
		mappedModel := raw.MapItemToModel(id, resolveData[id])
		result = append(result, *mappedModel)
	}

	fmt.Printf("ParsedVal: %+v \n", result[0])

	return result
}

func (bo *BondsOrder) UpdateItem () {}

func (bo *BondsOrder) MapItemToModel (id string, item map[string]string) *BondsOrder {
	height := item["order_height_" + id]
	price, priceErr := strconv.ParseInt(item["order_price_" + id], 10, 64)
	total, totalErr := strconv.ParseFloat(item["order_total_" + id], 64)
	filledtotal, filledTotalErr := strconv.ParseFloat(item["order_filled_total_" + id], 64)
	status := item["order_status_" + id]

	if priceErr != nil {
		price = 0
	}
	if totalErr != nil {
		total = 0
	}
	if filledTotalErr != nil {
		filledtotal = 0
	}

	wavesContractPower := float64(constants.WAVES_CONTRACT_POW)

	return &BondsOrder {
		Height: height,
		Price: int(price),
		Total: float64(total / wavesContractPower),
		Filledtotal: float64(filledtotal / wavesContractPower),
		Timestamp: 1111, // TODO
		Owner: item["order_owner_" + id],
		Resttotal: (total - filledtotal) / wavesContractPower,
		Status: status,
		Amount: total / (float64(price) * wavesContractPower / 100),
		Filledamount: filledtotal / (float64(price) * wavesContractPower / 100),
		Restamount: (total - filledtotal) / (float64(price) * wavesContractPower / 100),
		Pairname: "usdn-usdnb",
		Type: "buy",
	}
}
package entities;

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ventuary-lab/cache-updater/src/constants"
)

type BondsOrder struct {
	DAppEntity

	tableName struct{} `pg:"f_bonds_orders"`

	Order_id string `pg:",pk"` 
	Height, Owner, Status, Pairname, Type string
	Index *int
	Price int
	Timestamp int64
	Total, Filledamount, Filledtotal, Resttotal, Amount, Restamount float64
}

func (bo *BondsOrder) GetKeys(regex *string) []string {
	id := unwrapDefaultRegex(regex, "([A-Za-z0-9]{40,50})")

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

func (this *BondsOrder) UpdateAll (nodeData *map[string]string) []BondsOrder {
	ids := []string{}
	result := []BondsOrder{}
	regexKeys := this.GetKeys(nil)
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

			validKeys := this.GetKeys(&matchedAddress)

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

	return result
}

func (bo *BondsOrder) UpdateItem () {}

func (bo *BondsOrder) Includes (s *[]BondsOrder, e *BondsOrder) bool {
	for _, a := range *s {
		fmt.Printf("A: %+v; \nE: %+v \n", a, *e)
		if a == *e {
            return true
        }
    }
    return false
}

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

	var index *int = nil;
	orderbook := strings.Split(item["orderbook"], "_")
	for orderbookindex, orderbookpos := range orderbook {
		if orderbookpos == id {
			index = &orderbookindex
			break
		}
	}

	wavesContractPower := float64(constants.WAVES_CONTRACT_POW)

	return &BondsOrder {
		Order_id: id,
		Height: height,
		Price: int(price),
		Total: float64(total / wavesContractPower),
		Filledtotal: float64(filledtotal / wavesContractPower),
		Timestamp: 1111, // TODO
		Index: index,
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
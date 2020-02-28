package entities

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"github.com/ventuary-lab/cache-updater/src/constants"
)

const BONDS_ORDERS_NAME = "f_bonds_orders"

type BondsOrder struct {
	DAppEntity

	tableName struct{} `pg:"f_bonds_orders"`

	OrderId string `pg:"order_id,pk"`
	Owner, Status, Pairname, Type string
	OrderPrev *string `pg:"order_prev"`
	OrderNext *string `pg:"order_next"`
	IsFirst bool `pg:"is_first"`
	Index *int
	Price int
	Height, Timestamp uint64
	DebugROI uint64 `pg:"debug_roi"`
	DebugPrice uint64 `pg:"debug_price"`
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
		"debug_order_roi_" + id,
		"debug_order_currentPrice_" + id,
		"order_prev_" + id,
		"order_next_" + id,
		"order_first",
	}
}

func (bo *BondsOrder) UpdateAll (nodeData *map[string]string) []*BondsOrder {
	var ids []string
	regexKeys := bo.GetKeys(nil)
	heightKey := regexKeys[0]
	heightRegex, heightRegexErr := regexp.Compile(heightKey)
	var nodeKeys []string
	resolveData := make(map[string]map[string]string)

	for k, _ := range *nodeData {
		for _, regexKey := range regexKeys {
			compiledRegex := regexp.MustCompile(regexKey)

			if len(compiledRegex.FindSubmatch([]byte(k))) == 0 {
				continue
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
			validKeys := bo.GetKeys(&matchedAddress)

			for _, validKey := range validKeys {
				for _, k := range nodeKeys {
					if k == validKey {
						resolveData[matchedAddress][k] = (*nodeData)[k]
					}
				}
			}
		}
	}

	result := make([]*BondsOrder, len(ids))
	if heightRegexErr != nil {
		return result
	}

	// raw := BondsOrder{}

	for index, id := range ids {
		mappedModel := bo.MapItemToModel(id, resolveData[id])
		// result = append(result, mappedModel)
		result[index] = mappedModel
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
	height, _ := strconv.ParseInt(item["order_height_" + id], 10, 64)
	price, priceErr := strconv.ParseInt(item["order_price_" + id], 10, 64)
	total, totalErr := strconv.ParseFloat(item["order_total_" + id], 64)
	filledtotal, filledTotalErr := strconv.ParseFloat(item["order_filled_total_" + id], 64)
	status := item["order_status_" + id]
	orderROI, _ := strconv.ParseInt(item["debug_order_roi_" + id], 10, 64)
	orderPrice, _ := strconv.ParseInt(item["debug_order_currentPrice_" + id], 10, 64)
	rawOrderPrev := item["order_prev_" + id]
    rawOrderNext := item["order_next_" + id]
    var orderPrev, orderNext *string
	firstOrderId := item["order_first"]

	orderNext = nil
	if rawOrderNext != "" {
		orderNext = &rawOrderNext
	}
	orderPrev = nil
	if rawOrderPrev != "" {
		orderPrev = &rawOrderPrev
	}

	if priceErr != nil {
		price = 0
	}
	if totalErr != nil {
		total = 0
	}
	if filledTotalErr != nil {
		filledtotal = 0
	}

	var index *int = nil
	orderbook := strings.Split(item["orderbook"], "_")
	for orderbookindex, orderbookpos := range orderbook {
		if orderbookpos == id {
			index = &orderbookindex
			break
		}
	}

	wavesContractPower := float64(constants.WAVES_CONTRACT_POW)

	return &BondsOrder {
		OrderId: id,
		Height: uint64(height),
		Price: int(price),
		Total: total / wavesContractPower,
		Filledtotal: filledtotal / wavesContractPower,
		Index: index,
		Owner: item["order_owner_" + id],
		Resttotal: (total - filledtotal) / wavesContractPower,
		Status: status,
		Amount: total / (float64(price) * wavesContractPower / 100),
		Filledamount: filledtotal / (float64(price) * wavesContractPower / 100),
		Restamount: (total - filledtotal) / (float64(price) * wavesContractPower / 100),
		Pairname: "usd-nb_usd-n",
		Type: "buy",
		DebugROI: uint64(orderROI),
		DebugPrice: uint64(orderPrice),
		OrderNext: orderNext,
		OrderPrev: orderPrev,
		IsFirst: id == firstOrderId,
	}
}
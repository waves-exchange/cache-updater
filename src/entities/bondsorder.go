package entities

import (
	"fmt"
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

// Data Keys in blockchain
const (
	OrderHeightKey = "order_height_"
	OrderOwnerKey = "order_owner_"
	OrderPriceKey = "order_price_"
	OrderTotalKey = "order_total_"
	OrderFilledTotalKey = "order_filled_total_"
	OrderStatusKey = "order_status_"
	OrderBookKey = "orderbook"
	DebugOrderRoiKey = "debug_order_roi_"
	DebugOrderCurrentPriceKey = "debug_order_currentPrice_"
	OrderPrevKey = "order_prev_"
	OrderNextKey = "order_next_"
	OrderFirstKey = "order_first"
)

func (bo *BondsOrder) GetKeys(regex *string) []string {
	id := unwrapDefaultRegex(regex, "([A-Za-z0-9]{40,50})")

	return []string {
		OrderHeightKey + id,
		OrderOwnerKey + id,
		OrderPriceKey + id,
		OrderTotalKey + id,
		OrderFilledTotalKey + id,
		OrderStatusKey + id,
		OrderBookKey,
		DebugOrderRoiKey + id,
		DebugOrderCurrentPriceKey + id,
		OrderPrevKey + id,
		OrderNextKey + id,
		OrderFirstKey,
	}
}

func (bo *BondsOrder) UpdateAll (nodeData *map[string]string) []*BondsOrder {
	ids, resolveData, _ := UpdateAll(nodeData, bo.GetKeys)
	result := make([]*BondsOrder, len(ids))

	for index, id := range ids {
		mappedModel := bo.MapItemToModel(id, resolveData[id])
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

func (bo *BondsOrder) debugPriceDiff (price int64) float64 {
	return float64(price) / float64(constants.NEUTRINO_PRICE_DEC)
}

func (bo *BondsOrder) MapItemToModel (id string, item map[string]string) *BondsOrder {
	height, _ := strconv.ParseInt(item[OrderHeightKey + id], 10, 64)
	price, priceErr := strconv.ParseInt(item[OrderPriceKey + id], 10, 64)
	total, totalErr := strconv.ParseFloat(item[OrderTotalKey + id], 64)
	filledtotal, filledTotalErr := strconv.ParseFloat(item[OrderFilledTotalKey + id], 64)
	status := item[OrderStatusKey + id]
	orderROI, _ := strconv.ParseInt(item[DebugOrderRoiKey + id], 10, 64)
	orderPrice, _ := strconv.ParseInt(item[DebugOrderCurrentPriceKey + id], 10, 64)
	rawOrderPrev := item[OrderPrevKey + id]
    rawOrderNext := item[OrderNextKey + id]
    var orderPrev, orderNext *string
	firstOrderId := item[OrderFirstKey]

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
	orderbook := strings.Split(item[OrderBookKey], "_")
	for orderbookindex, orderbookpos := range orderbook {
		if orderbookpos == id {
			index = &orderbookindex
			break
		}
	}

	wavesContractPower := float64(constants.WAVES_CONTRACT_POW)

	bondsOrder := &BondsOrder {
		OrderId: id,
		Height: uint64(height),
		Price: int(price),
		Total: total / wavesContractPower,
		Index: index,
		Owner: item[OrderOwnerKey + id],
		Status: status,
		Resttotal: (total - filledtotal) / wavesContractPower,
		Filledtotal: filledtotal / wavesContractPower,
		Amount: total / (float64(price) * wavesContractPower / float64(constants.OLD_NEUTRINO_PRICE_DEC)),
		Filledamount: filledtotal / (float64(price) * wavesContractPower / float64(constants.OLD_NEUTRINO_PRICE_DEC)),
		Restamount: (total - filledtotal) / (float64(price) * wavesContractPower / float64(constants.OLD_NEUTRINO_PRICE_DEC)),
		Pairname: "usd-nb_usd-n",
		Type: "buy",
		DebugROI: uint64(orderROI),
		DebugPrice: uint64(orderPrice),
		OrderNext: orderNext,
		OrderPrev: orderPrev,
		IsFirst: id == firstOrderId,
	}

	if priceDiff := bo.debugPriceDiff(orderPrice); priceDiff > constants.OLD_PRICE_DECIMAL_DIFF {
		diffCoeff := float64(constants.NEUTRINO_PRICE_DEC / constants.OLD_NEUTRINO_PRICE_DEC)
		bondsOrder.Price /= int(diffCoeff)
		bondsOrder.Amount *= diffCoeff
		bondsOrder.Filledamount *= diffCoeff
		bondsOrder.Restamount *= diffCoeff
		bondsOrder.DebugPrice /= uint64(diffCoeff)

		return bondsOrder
	} else {
		return bondsOrder
	}
}
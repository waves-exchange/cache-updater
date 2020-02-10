package entities;

import (
	"strconv"
	// "github.com/ventuary-lab/cache-updater/enums"
)

type NeutrinoOrder struct {
	DAppEntity

	tableName struct{} `pg:"f_neutrino_orders"`

	Currency, Owner, status, Type, Order_id string
	Height uint64
	Ordernext, Orderprev *string
	Resttotal, Total int64
	// Status enums.OrderStatusEnum
	// Type enums.OrderTypeEnum
	Isfirst, Islast bool
}

func (no *NeutrinoOrder) GetKeys(regex *string) []string {
	id := unwrapDefaultRegex(regex, "([A-Za-z0-9]{40,50})")

	return []string {
		"order_height_" + id,
		"order_owner_" + id,
		"order_total_" + id,
		"order_filled_total_" + id,
		"order_status_" + id,
		"order_prev_" + id,
		"order_next_" + id,
		"orderbook",
		"order_first",
		"order_last",
	}
}

func (no *NeutrinoOrder) MapItemToModel (id string, item map[string]string) *NeutrinoOrder {
	height, _ := strconv.ParseInt(item["order_height_" + id], 10, 64)
	owner := item["order_owner_" + id]
	status := item["order_status_" + id]
	total, totalErr := strconv.ParseInt(item["order_total_" + id], 10, 64)
	filledtotal, filledtotalErr := strconv.ParseInt(item["order_filled_total_" + id], 10, 64)
	currency := "usd-nb"
	orderType := "liquidate"
	
	if totalErr != nil {
		total = 0
	}
	if filledtotalErr != nil {
		filledtotal = 0
	}

	rawOrderNext := item["order_next_" + id]
	rawOrderPrev := item["order_prev_" + id]
	var orderNext, orderPrev *string
	if rawOrderNext == "" { orderNext = nil } else { orderNext = &rawOrderNext }
	if rawOrderPrev == "" { orderPrev = nil } else { orderPrev = &rawOrderPrev }

	return &NeutrinoOrder{
		Height: uint64(height),
		Currency: currency,
		Owner: owner,
		Total: total,
		status: status,
		Resttotal: total - filledtotal,
		Type: orderType,
		Ordernext: orderNext,
		Orderprev: orderPrev,
		Isfirst: id == item["order_first"],
		Islast: id == item["order_last"],
	}
}
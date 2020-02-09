package entities;

import (
	"github.com/ventuary-lab/cache-updater/enums"
)

type NeutrinoOrder struct {
	DAppEntity

	tableName struct{} `pg:"f_neutrino_orders"`

	Currency, Owner, Status, Type, Order_id string
	Height uint64
	Ordernext, Orderprev *string
	Resttotal, Total int64
	// Status enums.OrderStatusEnum
	// Type enums.OrderTypeEnum
	Isfirst, Islast bool
}

func (this *NeutrinoOrder) GetKeys(regex *string) []string {
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

func (this *NeutrinoOrder) MapItemToModel (id string, item map[string]string) *NeutrinoOrder {
	height, _ := strconv.ParseInt(item["order_height_" + id], 10, 64)
	owner := item["order_owner_" + id]
	status := item["order_status_" + id]
	orderNext := item["order_next_" + id]
	orderPrev := item["order_prev_" + id]
	total, totalErr := strconv.ParseInt(item["order_total_" + id], 10, 64)
	filledtotal, filledTotalErr := strconv.ParseInt(item["order_filled_total_" + id, 10, 64)
	currency := "usd-nb"
	orderType := "liquidate"

	if totalErr != nil {
		total = 0
	}
	if filledTotalErr != nil {
		filledtotal = 0
	}
	if orderNext == "" {
		orderNext = nil
	}
	if orderPrev == "" {
		orderPrev = nil
	}

	return &NeutrinoOrder{
		Height: height,
		Currency: currency,
		Owner: owner,
		Total: total,
		Resttotal: total - filledTotal,
		Type: orderType,
		OrderNext: orderNext,
		OrderPrev: orderPrev,
		IsFirst: id == item["order_first"],
		IsLast: id == item["order_last"]
	}
}
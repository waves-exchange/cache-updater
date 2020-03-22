package entities

import (
	"strconv"
	// "github.com/ventuary-lab/cache-updater/enums"
)

const NEUTRINO_ORDERS_NAME = "f_neutrino_orders"

type NeutrinoOrder struct {
	tableName struct{} `pg:"f_neutrino_orders"`

	OrderId *string `pg:"order_id,pk"`
	Currency, Owner, Status, Type *string
	Height uint64
	OrderPrev *string `pg:"order_prev"`
	OrderNext *string `pg:"order_next"`
	RestTotal int64 `pg:"resttotal"`
	Total int64 `pg:"total"`
	// Status enums.OrderStatusEnum
	// Type enums.OrderTypeEnum

	IsFirst bool `pg:"is_first"`
	IsLast bool `pg:"is_last"`
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

func (no *NeutrinoOrder) UpdateAll (nodeData *map[string]string) []*NeutrinoOrder {
	ids, resolveData, _ := UpdateAll(nodeData, no.GetKeys)
	result := make([]*NeutrinoOrder, len(ids))

	for index, id := range ids {
		mappedModel := no.MapItemToModel(id, resolveData[id])
		result[index] = mappedModel
	}

	return result
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
		Currency: &currency,
		Owner: &owner,
		Total: total,
		Status: &status,
		RestTotal: total - filledtotal,
		Type: &orderType,
		OrderNext: orderNext,
		OrderPrev: orderPrev,
		IsFirst: id == item["order_first"],
		IsLast: id == item["order_last"],
	}
}
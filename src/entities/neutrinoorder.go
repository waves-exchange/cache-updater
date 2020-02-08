package entities;

import (
	"github.com/ventuary-lab/cache-updater/enums"
)

type NeutrinoOrder struct {
	DAppEntity

	tableName struct{} `pg:"f_neutrino_orders"`

	Height, Currency, Owner, Total, Order_id string
	Ordernext, Orderprev *string
	Timestamp int64
	Status enums.OrderStatusEnum
	Resttotal int
	Type enums.OrderTypeEnum
	Isfirst, Islast bool
}

func (no NeutrinoOrder) GetKeys(regex *string) []string {
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


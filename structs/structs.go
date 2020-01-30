package structs;

import (
	"github.com/ventuary-lab/cache-updater/enums"
)

type BondsOrder struct {
	Height, Owner, Status, Index, Pairname, Type string
	Price int
	Timestamp int64
	Total, Filledamount, Filledtotal, Resttotal, Amount, Restamount float64
}

type NeutrinoOrder struct {
	Height, Currency, Owner, Total, Order_id string
	Ordernext, Orderprev *string
	Timestamp int64
	Status enums.OrderStatusEnum
	Resttotal int
	Type enums.OrderTypeEnum
	Isfirst, Islast bool
}

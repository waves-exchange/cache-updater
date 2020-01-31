package entities;

import (
	"strconv"
)

type BondsOrder struct {
	Height, Owner, Status, Index, Pairname, Type string
	Price int
	Timestamp int64
	Total, Filledamount, Filledtotal, Resttotal, Amount, Restamount float64
}

func (bo BondsOrder) GetKeys(id string) []string {
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


func (bo *BondsOrder) MapItemToModel (id string, item map[string]string) *BondsOrder {
	height := item["order_height_" + id];
	price, priceErr := strconv.ParseInt(item["order_price_" + id], 10, 64);
	total, totalErr := strconv.ParseFloat(item["order_total_" + id], 64);
	filledtotal, filledTotalErr := strconv.ParseFloat(item["order_filledtotal_" + id], 64);

	if totalErr == nil || filledTotalErr == nil || priceErr == nil {
		return nil;
	}

	return &BondsOrder {
		Height: height,
		Price: int(price),
		Total: float64(total),
		Filledtotal: float64(filledtotal),
		Timestamp: 1111, // TODO
		Owner: item["order_owner_" + id],
		Resttotal: -10,
		Status: "new",
		Amount: 100,
		Filledamount: 101,
		Restamount: 102,
		Pairname: "usdn-usdnb",
		Type: "buy",
	}

    // async _prepareItem(id, item) {
    //     const index = item.orderbook
    //         .split('_')
    //         .filter(Boolean)
    //         .indexOf(id);

    //     const height = item['order_height_' + id];
    //     const price = item['order_price_' + id] || 0;
    //     const total = item['order_total_' + id] || 0;
    //     const filledTotal = item['order_filled_total_' + id] || 0;
    //     return {
    //         -height,
    //         -timestamp: (await this.heightListener.getTimestamps([height]))[height],
    //         -owner: item['order_owner_' + id],
    //         -price: Number(price),
    //         -total: _round(total / CurrencyEnum.getContractPow(CurrencyEnum.WAVES), 2),
    //         -filledTotal: _round(filledTotal / CurrencyEnum.getContractPow(CurrencyEnum.WAVES), 2),
    //         restTotal: _round((total - filledTotal) / CurrencyEnum.getContractPow(CurrencyEnum.WAVES), 2),
    //         status: item['order_status_' + id],
    //         index: index !== -1 ? index : null,
    //         amount: _round(total / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100)), // Bonds amount
    //         filledAmount: _round(filledTotal / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100), 2),
    //         restAmount: _round((total - filledTotal) / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100), 2),
    //         pairName: this.pairName,
    //         type: OrderTypeEnum.BUY,
    //     };
	// }
	
	// item: {                                                                                             r
	// 	order_height_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ: '1868066',
	// 	order_owner_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ: '3PGmja5rWBPiQ7n9eLSBgQBd6EzTmUFgddB',
	// 	order_price_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ: '55',
	// 	order_total_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ: '1434000000',
	// 	order_status_zyyXjxzKajKht1wJUGmBrWPcrjRbVbgFRtjQHJEHymJ: 'canceled',
	// 	orderbook: ''
	//   } 
}
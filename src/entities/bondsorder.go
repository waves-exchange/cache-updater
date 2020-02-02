package entities;

import (
	"strconv"
	"math"
	"regexp"
	"fmt"
	"github.com/ventuary-lab/cache-updater/src/constants"
)

type BondsOrder struct {
	Height, Owner, Status, Index, Pairname, Type string
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

func (bo *BondsOrder) UpdateAll (nodeData *map[string]string) {

	ids := []string{}
	defaultRawRegex := "([A-Za-z0-9]{40,50})"
	heightKey := bo.GetKeys(defaultRawRegex)[0]
	heightRegex, heightRegexErr := regexp.Compile(heightKey)
	nodeKeys := []string{}
	// resolveData := map[string] {}
	resolveData := make(map[string](map[string]string))

	for k, nodeVal := range *nodeData {
		resolveData[k] = map[string]string{}
		nodeKeys = append(nodeKeys, k)

		heightRegexSubmatches := heightRegex.FindSubmatch([]byte(k))

		if len(heightRegexSubmatches) < 2 {
			continue
		}

		matchedAddress := string(heightRegexSubmatches[1])

		if matchedAddress != "" {
			ids = append(ids, matchedAddress)

			validKeys := bo.GetKeys(matchedAddress)

			for _, validKey := range validKeys {
				if StrArrayContains(nodeKeys, validKey) {
					resolveData[validKey][k] = nodeVal
				}
			}
		}
	}

	if heightRegexErr != nil {
		return
	}

	fmt.Printf("HeightRegex: %v \n", heightRegex)
	fmt.Printf("ID #1: %v; COUNT: %v \n", ids[0], len(ids))

}

func (bo *BondsOrder) UpdateItem () {}

func (bo *BondsOrder) MapItemToModel (id string, item map[string]string) *BondsOrder {
	height := item["order_height_" + id]
	price, priceErr := strconv.ParseInt(item["order_price_" + id], 10, 64)
	total, totalErr := strconv.ParseFloat(item["order_total_" + id], 64)
	filledtotal, filledTotalErr := strconv.ParseFloat(item["order_filledtotal_" + id], 64)
	status := item["order_status_" + id]

	if priceErr == nil {
		price = 0
	}
	if totalErr == nil {
		total = 0
	}
	if filledTotalErr == nil {
		filledtotal = 0
	}

	// func ComputeTotal(t, p )N

	wavesContractPower := constants.WAVES_CONTRACT_POW

	resttotal := math.Round((total - filledtotal) / float64(wavesContractPower))
	// restTotal: _round((total - filledTotal) / CurrencyEnum.getContractPow(CurrencyEnum.WAVES), 2),
	amount := math.Round(total / (float64(price) * float64(wavesContractPower) / 100))
	// amount: _round(total / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100)), // Bonds amount
	filledAmount := math.Round(filledtotal / (float64(price) * float64(wavesContractPower) / 100))
	// filledAmount: _round(filledTotal / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100), 2),
	restAmount := math.Round((total - filledtotal) / (float64(price) * float64(wavesContractPower) / 100))
	// restAmount: _round((total - filledTotal) / (price * CurrencyEnum.getContractPow(CurrencyEnum.WAVES) / 100), 2),

	return &BondsOrder {
		Height: height,
		Price: int(price),
		Total: float64(total),
		Filledtotal: float64(filledtotal),
		Timestamp: 1111, // TODO
		Owner: item["order_owner_" + id],
		Resttotal: resttotal,
		Status: status,
		Amount: amount,
		Filledamount: filledAmount,
		Restamount: restAmount,
		Pairname: "usdn-usdnb",
		Type: "buy",
	}

    // async _prepareItem(id, item) {
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
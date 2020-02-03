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
	regexKeys := bo.GetKeys(defaultRawRegex)
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

			validKeys := bo.GetKeys(matchedAddress)

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
		return
	}

	
	raw := BondsOrder{}
	fmt.Printf("ID: %v; Val: %v \n", ids[0])
	fmt.Printf("ParsedVal: %+v \n", raw.MapItemToModel(ids[0], resolveData[ids[0]]))

	/*
	{
		height: '1763038',
		timestamp: 1571838501734,
		owner: '3PHNudHo6zeuFgop35TVAtnVHitx8SGPRiC',
		price: 70,
		total: 0,
		filledTotal: 0,
		restTotal: 0,
		status: 'filled',
		index: null,
		amount: 0,
		filledAmount: 0,
		restAmount: 0,
		pairName: 'usd-nb_usd-n',
		type: 'buy'
	}
	*/

}

func (bo *BondsOrder) UpdateItem () {}

func (bo *BondsOrder) MapItemToModel (id string, item map[string]string) *BondsOrder {
	height := item["order_height_" + id]
	price, priceErr := strconv.ParseInt(item["order_price_" + id], 10, 64)
	total, totalErr := strconv.ParseFloat(item["order_total_" + id], 64)
	filledtotal, filledTotalErr := strconv.ParseFloat(item["order_filledtotal_" + id], 64)
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

	fmt.Printf("Price: %v, Total: %v, FilledTotal: %v \n", price, total, filledtotal)
	// fmt.Printf("Price: %v, Total: %v, FilledTotal: %v \n", price, total, filledtotal)

	// func ComputeTotal(t, p )N

	wavesContractPower := constants.WAVES_CONTRACT_POW

	total = math.Round(total / float64(wavesContractPower))
	// total: _round(total / CurrencyEnum.getContractPow(CurrencyEnum.WAVES), 2),
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
}
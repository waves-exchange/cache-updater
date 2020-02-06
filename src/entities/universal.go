package entities;

type DAappNumberRecord struct {
	Key, Type string
	Value *int
}

type DAppStringRecord struct {
	Key, Type string
	Value *string
}

func StrArrayContains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func unwrapDefaultRegex (rawregex *string, defaultRegex string) string {
	if rawregex == nil {
		return defaultRegex
	} else {
		return *rawregex
    }
}

// async _prepareItem(id, item) {
//     const orderNext = item['order_next_' + id] || null;
//     const orderPrev = item['order_prev_' + id] || null;
//     const height = item['order_height_' + id];
//     const total = item['order_total_' + id] || 0;
//     const filledTotal = item['order_filled_total_' + id] || 0;
//     return {
//         height,
//         currency: this.pairName.split('_')[0],
//         timestamp: (await this.heightListener.getTimestamps([height]))[height],
//         owner: item['order_owner_' + id],
//         status: item['order_status_' + id],
//         total,
//         restTotal: total - filledTotal,
//         type: OrderTypeEnum.LIQUIDATE,
//         orderNext,
//         orderPrev,
//         isFirst: id == item.order_first,
//         isLast: id == item.order_last,
//     };
// }
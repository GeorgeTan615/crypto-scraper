package types

var StringToExchange = map[string]Exchange{
	"LUNO":    LUNO,
	"BINANCE": BINANCE,
}

var StringToType = map[string]Type{
	"ORDER_BOOK": ORDER_BOOK,
}

func MapExchanges(exchanges []string) []Exchange {
	res := make([]Exchange, 0, len(exchanges))

	for _, exchange := range exchanges {
		e, ok := StringToExchange[exchange]

		if !ok {
			continue
		}

		res = append(res, e)
	}

	return res
}

func MapTypes(types []string) []Type {
	res := make([]Type, 0, len(types))

	for _, typ := range types {
		t, ok := StringToType[typ]

		if !ok {
			continue
		}

		res = append(res, t)
	}

	return res
}

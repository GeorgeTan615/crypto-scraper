package bybit

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/reconquest/karma-go"
)

type orderBookResponse struct {
	Timestamp   time.Time
	BidPrice    string
	BidQuantity string
	AskPrice    string
	AskQuantity string
}

type bookPriceLevel struct {
	Price  string
	Amount string
}

func (r *bookPriceLevel) UnmarshalJSON(data []byte) error {
	var msg []string

	if err := json.Unmarshal(data, &msg); err != nil {
		return karma.Describe("body", string(data)).Format(err, "unmarshal order book response")
	}

	r.Price = msg[0]
	r.Amount = msg[1]

	return nil
}

func (r *orderBookResponse) GetHeaders() []string {
	return []string{
		"timestamp",
		"bidPrice",
		"bidQuantity",
		"askPrice",
		"askQuantity",
	}
}

func (r *orderBookResponse) Get() []string {
	return []string{
		r.Timestamp.Format(time.RFC3339Nano),
		r.BidPrice,
		r.BidQuantity,
		r.AskPrice,
		r.AskQuantity,
	}
}

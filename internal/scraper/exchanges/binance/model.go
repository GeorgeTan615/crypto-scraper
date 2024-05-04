package binance

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

func (r *orderBookResponse) UnmarshalJSON(data []byte) error {
	var msg struct {
		BidPrice    string `json:"bidPrice"`
		BidQuantity string `json:"bidQty"`
		AskPrice    string `json:"askPrice"`
		AskQuantity string `json:"askQty"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return karma.Describe("body", string(data)).Format(err, "unmarshal order book response")
	}

	r.BidPrice = msg.BidPrice
	r.BidQuantity = msg.BidQuantity
	r.AskPrice = msg.AskPrice
	r.AskQuantity = msg.AskQuantity
	r.Timestamp = time.Now()

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

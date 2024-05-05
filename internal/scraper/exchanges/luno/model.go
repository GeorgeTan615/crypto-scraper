package luno

import (
	"encoding/json"
	"fmt"
	"time"

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
	type order struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	}

	var msg struct {
		// Timestamp int64   `json:"timestamp"`
		Asks []order `json:"asks"`
		Bids []order `json:"bids"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return karma.Describe("body", string(data)).Format(err, "unmarshal order book response")
	}

	if len(msg.Bids) == 0 {
		return fmt.Errorf("expected length of bids at least 1, got: %v", msg.Bids)
	}

	if len(msg.Asks) == 0 {
		return fmt.Errorf("expected length of asks at least 1, got: %v", msg.Asks)
	}

	r.BidPrice = msg.Bids[0].Price
	r.BidQuantity = msg.Bids[0].Volume
	r.AskPrice = msg.Asks[0].Price
	r.AskQuantity = msg.Asks[0].Volume
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

package bybit

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/crypto-scraper/internal/scraper/message"
	"github.com/crypto-scraper/internal/scraper/web"
	"github.com/goccy/go-json"
	"github.com/reconquest/karma-go"
)

const _orderBookAPIPath = "/v5/market/orderbook"

type BybitOrderBookScrapper struct {
	client *web.Client
}

func NewBybitOrderBookScrapper() *BybitOrderBookScrapper {
	return &BybitOrderBookScrapper{
		client: web.NewClient(),
	}
}

func (s *BybitOrderBookScrapper) Scrape(
	ctx context.Context,
	req *message.ScrapeRequest,
) (message.CSVData, error) {
	url, err := s.prepareURL(req.Symbol)
	if err != nil {
		return nil, karma.Format(err, "prepare url")
	}

	apiReq, err := web.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, karma.Format(err, "create request")
	}

	resp, err := s.client.Do(apiReq)
	if err != nil {
		return nil, karma.Format(err, "make request")
	}
	defer resp.Close()

	bodyBytes, err := resp.BodyBytes()
	if err != nil {
		return nil, karma.Format(err, "read body bytes")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, karma.Describe("body", string(bodyBytes)).Reason("response not ok")
	}

	var data struct {
		RetMsg string `json:"retMsg"`
		Result struct {
			Asks []bookPriceLevel `json:"a"`
			Bids []bookPriceLevel `json:"b"`
		} `json:"result"`
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, karma.Format(err, "unmarshal body bytes")
	}

	if data.RetMsg != "OK" {
		return nil, karma.Describe("body", string(bodyBytes)).Reason("data is not OK")
	}

	if len(data.Result.Asks) == 0 {
		return nil, karma.Describe("body", string(bodyBytes)).Reason("asks is empty")
	}

	if len(data.Result.Bids) == 0 {
		return nil, karma.Describe("body", string(bodyBytes)).Reason("bids is empty")
	}

	return &orderBookResponse{
		Timestamp:   time.Now(),
		BidPrice:    data.Result.Bids[0].Price,
		BidQuantity: data.Result.Bids[0].Amount,
		AskPrice:    data.Result.Asks[0].Price,
		AskQuantity: data.Result.Asks[0].Amount,
	}, nil
}

func (s *BybitOrderBookScrapper) prepareURL(symbol string) (string, error) {
	parsedURL, err := url.Parse(_bybitBaseURL + _orderBookAPIPath)
	if err != nil {
		return "", karma.Format(err, "parse url %s", _bybitBaseURL+_orderBookAPIPath)
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	query.Set("category", "spot")
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

package binance

import (
	"context"
	"net/http"
	"net/url"

	"github.com/crypto-scraper/internal/scraper/message"
	"github.com/crypto-scraper/internal/scraper/web"
	"github.com/goccy/go-json"
	"github.com/reconquest/karma-go"
)

const _orderBookAPIPath = "/api/v3/ticker/bookTicker"

type BinanceOrderBookScrapper struct {
	client *web.Client
}

func NewBinanceOrderBookScrapper() *BinanceOrderBookScrapper {
	return &BinanceOrderBookScrapper{
		client: web.NewClient(),
	}
}

func (s *BinanceOrderBookScrapper) Scrape(
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

	var data orderBookResponse

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, karma.Format(err, "unmarshal body bytes")
	}

	return &data, nil
}

func (s *BinanceOrderBookScrapper) prepareURL(symbol string) (string, error) {
	parsedURL, err := url.Parse(_binanceBaseURL + _orderBookAPIPath)
	if err != nil {
		return "", karma.Format(err, "parse url %s", _binanceBaseURL+_orderBookAPIPath)
	}

	query := parsedURL.Query()
	query.Set("symbol", symbol)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

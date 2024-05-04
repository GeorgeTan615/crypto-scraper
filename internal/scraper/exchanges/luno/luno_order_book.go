package luno

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/crypto-scraper/internal/scraper/message"
	"github.com/crypto-scraper/internal/scraper/web"
	"github.com/reconquest/karma-go"
)

const _orderBookAPIPath = "/api/1/orderbook_top"

type LunoOrderBookScrapper struct {
	client               *web.Client
	symbolReplacementMap map[string]string
}

func NewLunoOrderBookScrapper() *LunoOrderBookScrapper {
	return &LunoOrderBookScrapper{
		client: web.NewClient(),
		symbolReplacementMap: map[string]string{
			"BTC": "XBT",
		},
	}
}

func (s *LunoOrderBookScrapper) Scrape(
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

func (s *LunoOrderBookScrapper) prepareURL(symbol string) (string, error) {
	parsedURL, err := url.Parse(_lunoBaseURL + _orderBookAPIPath)
	if err != nil {
		return "", karma.Format(err, "parse url %s", _lunoBaseURL+_orderBookAPIPath)
	}

	for initialSymbol, newSymbol := range s.symbolReplacementMap {
		if strings.Contains(symbol, initialSymbol) {
			symbol = strings.ReplaceAll(symbol, initialSymbol, newSymbol)
			break
		}
	}

	query := parsedURL.Query()
	query.Set("pair", symbol)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

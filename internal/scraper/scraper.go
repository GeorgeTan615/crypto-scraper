package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/crypto-scraper/internal/scheduler"
	"github.com/crypto-scraper/internal/scraper/exchanges/binance"
	"github.com/crypto-scraper/internal/scraper/exchanges/bybit"
	"github.com/crypto-scraper/internal/scraper/exchanges/luno"
	"github.com/crypto-scraper/internal/scraper/message"
	"github.com/crypto-scraper/internal/types"
	"github.com/reconquest/karma-go"
)

const _scraperTimeout = 10 * time.Second

type (
	Scrapper interface {
		Scrape(ctx context.Context, req *message.ScrapeRequest) (message.CSVData, error)
	}

	ScrapperManager struct {
		registry map[types.Exchange]map[types.Type]Scrapper
	}
)

func NewScrapperManager() *ScrapperManager {
	return &ScrapperManager{
		registry: map[types.Exchange]map[types.Type]Scrapper{
			types.BINANCE: {
				types.ORDER_BOOK: binance.NewBinanceOrderBookScrapper(),
			},
			types.LUNO: {
				types.ORDER_BOOK: luno.NewLunoOrderBookScrapper(),
			},
			types.BYBIT: {
				types.ORDER_BOOK: bybit.NewBybitOrderBookScrapper(),
			},
		},
	}
}

func (sm *ScrapperManager) Start(
	ctx context.Context,
	exchangesConfig map[types.Exchange]map[types.Type]map[string]time.Duration,
) {
	var wg sync.WaitGroup

	for exchange, typsMap := range exchangesConfig {
		for typ, symbolsMap := range typsMap {
			for symbol, interval := range symbolsMap {
				logContext := karma.
					Describe("exchange", exchange).
					Describe("type", typ)

				t, ok := sm.registry[exchange]
				if !ok {
					log.Println(karma.Format(logContext, "exchange not registered"))
					continue
				}

				s, ok := t[typ]
				if !ok {
					log.Println(karma.Format(logContext, "type not registered"))
					continue
				}

				wg.Add(1)
				go func(exchange types.Exchange, typ types.Type, symbol string, interval time.Duration) {
					defer wg.Done()

					sm.StartScrapper(
						ctx,
						s,
						exchange,
						typ,
						symbol,
						interval,
					)
				}(exchange, typ, symbol, interval)
			}
		}
	}

	wg.Wait()
}

func (sm *ScrapperManager) StartScrapper(
	ctx context.Context,
	scraper Scrapper,
	exchange types.Exchange,
	typ types.Type,
	symbol string,
	interval time.Duration,
) {
	logContext := karma.
		Describe("exchange", exchange).
		Describe("type", typ).
		Describe("symbol", symbol).
		Describe("interval", interval)

	// create destinationn file
	csvFile, err := os.Create(
		fmt.Sprintf(
			"%s_%s_%s.csv",
			exchange,
			typ,
			time.Now().Format("2006-01-02_15_04_05")))
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return
	}
	defer csvFile.Close()

	var (
		csvWriter         = csv.NewWriter(csvFile)
		areHeadersWritten = false
		chRes             = make(chan []string, 1000)
		wg                sync.WaitGroup
	)

	// collect data and write to csv
	wg.Add(1)
	go func() {
		defer wg.Done()

		for data := range chRes {
			err := csvWriter.Write(data)
			if err != nil {
				log.Println(logContext.
					Describe("operation", "write csv").
					Describe("data", data).
					Reason(err))
				continue
			}

			csvWriter.Flush()
		}
	}()

	// start scheduler to scrape data periodically
	scheduler := scheduler.NewScheduler(interval, func(c context.Context) {
		ctx, cancel := context.WithTimeout(c, _scraperTimeout)
		defer cancel()

		csvData, err := scraper.Scrape(ctx, &message.ScrapeRequest{Symbol: symbol})
		if err != nil {
			log.Println(logContext.Describe("operation", "scrape").Reason(err))
			return
		}

		// write headers if file is empty
		if !areHeadersWritten {
			chRes <- csvData.GetHeaders()
			areHeadersWritten = true
		}

		// write scraped data to csv
		chRes <- csvData.Get()
	})

	log.Println(logContext.Reason("start scraping"))
	scheduler.Start(ctx)
	close(chRes)
	log.Println(logContext.Reason("stop scraping"))
}

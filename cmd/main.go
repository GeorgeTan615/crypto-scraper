package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/crypto-scraper/internal/scraper"
	"github.com/crypto-scraper/internal/types"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
)

var usage = `crypto_scraper

Scrapes data off exchanges.

Usage:
  crypto_scraper --exchanges <exchanges> --types <types> --interval <interval> --symbols <symbols>
  crypto_scraper -h | --help
  crypto_scraper --version

Required options:
  -e --exchanges <exchanges>	   	Exchanges to scrape from (delimited by ',')
  -t --types <types>						Data types to scrape from exchanges (delimited by ',')
  -s --symbols <symbols>				Symbols to scrape (delimited by ',')
  -i --interval <interval>  			Scrape interval (seconds)

Options:
  -h --help     Show this screen.
  --version     Show version.`

type Config struct {
	Exchanges []types.Exchange
	Types     []types.Type
	Symbols   []string
	Interval  time.Duration
}

func main() {
	config, err := parseCLIArgs()
	if err != nil {
		log.Fatalln(karma.Format(err, "parse cli args"))
	}

	sm := scraper.NewScrapperManager()
	sm.Init()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)
	defer cancel()

	sm.Start(ctx, config.Exchanges, config.Types, config.Symbols, config.Interval)
}

func parseCLIArgs() (*Config, error) {
	args, err := docopt.ParseArgs(usage, nil, "")
	if err != nil {
		return nil, karma.Format(err, "parse arguments")
	}

	var config struct {
		Exchanges string `docopt:"--exchanges"`
		Types     string `docopt:"--types"`
		Interval  string `docopt:"--interval"`
		Symbols   string `docopt:"--symbols"`
	}

	err = args.Bind(&config)
	if err != nil {
		return nil, karma.Format(err, "bind arguments")
	}

	exchangesStr := strings.Split(config.Exchanges, ",")
	typeStr := strings.Split(config.Types, ",")
	symbols := strings.Split(config.Symbols, ",")

	exchanges := types.MapExchanges(exchangesStr)
	if len(exchanges) == 0 {
		return nil, fmt.Errorf("no allowed exchanges: %v", config.Exchanges)
	}

	types := types.MapTypes(typeStr)
	if len(types) == 0 {
		return nil, fmt.Errorf("no allowed types: %v", config.Types)
	}

	interval, err := strconv.Atoi(config.Interval)
	if err != nil {
		return nil, karma.Format(err, "parse int: %v", config.Interval)
	}

	return &Config{
		Exchanges: exchanges,
		Types:     types,
		Symbols:   symbols,
		Interval:  time.Second * time.Duration(interval),
	}, nil
}

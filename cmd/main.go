package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crypto-scraper/internal/scraper"
	"github.com/crypto-scraper/internal/types"
	"gopkg.in/yaml.v3"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
)

var usage = `crypto_scraper

Scrapes data off exchanges.

Usage:
  crypto_scraper --config <config>
  crypto_scraper -h | --help
  crypto_scraper --version

Required options:
  -c --config <config>					File path to config file

Options:
  -h --help     Show this screen.
  --version     Show version.`

type config struct {
	ExchangesScrapeConfig map[types.Exchange]map[types.Type]map[string]time.Duration `yaml:"exchanges"`
}

func main() {
	config, err := parseCLIArgs()
	if err != nil {
		log.Fatalln(karma.Format(err, "parse cli args"))
		return
	}

	sm := scraper.NewScrapperManager()
	sm.Init()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)
	defer cancel()

	sm.Start(ctx, config.ExchangesScrapeConfig)
}

func parseCLIArgs() (*config, error) {
	args, err := docopt.ParseArgs(usage, nil, "")
	if err != nil {
		return nil, karma.Format(err, "parse arguments")
	}

	var arguments struct {
		ConfigFilePath string `docopt:"--config"`
	}

	err = args.Bind(&arguments)
	if err != nil {
		return nil, karma.Format(err, "bind arguments")
	}

	yamlFile, err := os.ReadFile(arguments.ConfigFilePath)
	if err != nil {
		return nil, karma.Format(err, "read config file")
	}

	var config config

	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, karma.Format(err, "unmarshal yaml file")
	}

	return &config, nil
}

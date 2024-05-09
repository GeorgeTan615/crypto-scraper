# Crypto Scraper
A scraper catered to scrape crypto data from data sources.

Data would be stored into a `.csv` file and data points could then be utilised for backtesting/constructing algorithmic trading strategies.

The scraper is intentionally structured in a way that it's flexible enough to cater for various combinations of exchanges, data types, and intervals. Check out `config.yaml` for an example configuration.

## How to Run
`go run ./cmd/main.go -c config.yaml`
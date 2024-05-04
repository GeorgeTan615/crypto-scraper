package message

type ScrapeRequest struct {
	Symbol string
}

type CSVData interface {
	GetHeaders() []string
	Get() []string
}

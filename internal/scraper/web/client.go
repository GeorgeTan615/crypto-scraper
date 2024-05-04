package web

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/reconquest/karma-go"
)

func NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	userAgent := browser.Random()
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept-Encoding", "gzip")

	return req, nil
}

type Response struct {
	*http.Response
}

func (r *Response) BodyBytes() ([]byte, error) {
	reader := r.Body

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, karma.Format(err, "create gzip reader")
		}
		reader = gzipReader
	}

	return io.ReadAll(reader)
}

func (r *Response) Close() {
	if r.Body != nil {
		r.Body.Close()
	}
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) Do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: resp,
	}, nil
}

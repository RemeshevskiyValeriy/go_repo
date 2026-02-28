package httpx

import (
	"context"
	"net/http"
	"time"

	"example.com/tech-ip-sem2/shared/middleware"
)

type Client struct {
	httpClient *http.Client
}

func New(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// прокидываем request-id
	if rid := middleware.GetRequestID(ctx); rid != "" {
		req.Header.Set(middleware.HeaderRequestID, rid)
	}

	req = req.WithContext(ctx)

	return c.httpClient.Do(req)
}
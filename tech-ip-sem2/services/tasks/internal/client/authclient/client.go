package authclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"example.com/tech-ip-sem2/shared/httpx"
)

var ErrUnauthorized = errors.New("unauthorized")

type Client struct {
	baseURL string
	client  *httpx.Client
}

func New(baseURL string, client *httpx.Client) *Client {
	return &Client{
		baseURL: baseURL,
		client:  client,
	}
}

func (c *Client) Verify(ctx context.Context, authHeader string) error {
	req, err := http.NewRequest(
		http.MethodGet,
		c.baseURL+"/v1/auth/verify",
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := c.client.Do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("auth service error")
	}

	var body struct {
		Valid bool `json:"valid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	if !body.Valid {
		return ErrUnauthorized
	}

	return nil
}
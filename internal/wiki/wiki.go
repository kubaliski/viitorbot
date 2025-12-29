package wiki

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client defines the behavior used by handlers; allows easier testing.
type Client interface {
	GetRandomArticle() (string, string, string, error)
}

type client struct{}

// Endpoint can be overridden in tests.
var Endpoint = "https://es.wikipedia.org/api/rest_v1/page/random/summary"

func NewClient() Client { return &client{} }

type wikipediaResponse struct {
	Title       string `json:"title"`
	Extract     string `json:"extract"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

func (c *client) GetRandomArticle() (string, string, string, error) {
	url := Endpoint

	// Use an explicit client and set a User-Agent; Wikipedia may reject requests
	// without a sensible User-Agent header (403). Also set a timeout.
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("User-Agent", "viitorbot/1.0 (+https://example.org)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("read body error: %w", err)
	}

	var w wikipediaResponse
	if err := json.Unmarshal(body, &w); err != nil {
		return "", "", "", fmt.Errorf("json unmarshal error: %w", err)
	}

	return w.Title, w.Extract, w.ContentURLs.Desktop.Page, nil
}

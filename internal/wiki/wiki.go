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
	GetRandomArticleByDate(t time.Time) (string, string, string, *Evidence, error)
}

type client struct{}

// Endpoint can be overridden in tests.
var Endpoint = "https://es.wikipedia.org/api/rest_v1/page/random/summary"

// OnThisDayEndpoint formats with month, day (MM, DD)
var OnThisDayEndpoint = "https://es.wikipedia.org/api/rest_v1/feed/onthisday/all/%02d/%02d"

func NewClient() Client { return &client{} }

// Evidence represents provenance information for an on-this-day item.
type Evidence struct {
	Type string `json:"type"`
	Year int    `json:"year"`
	Text string `json:"text"`
}

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

// GetRandomArticleByDate queries the onthisday feed for the given date (month/day)
// and returns a random page from the aggregated lists (events/births/deaths/holidays).
func (c *client) GetRandomArticleByDate(t time.Time) (string, string, string, *Evidence, error) {
	url := fmt.Sprintf(OnThisDayEndpoint, int(t.Month()), t.Day())

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("User-Agent", "viitorbot/1.0 (+https://example.org)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", nil, fmt.Errorf("read body error: %w", err)
	}

	// Minimal types for onthisday feed
	type page struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
	}
	type item struct {
		Type  string `json:"type"`
		Year  int    `json:"year"`
		Text  string `json:"text"`
		Pages []page `json:"pages"`
	}
	var feed struct {
		Events   []item `json:"events"`
		Births   []item `json:"births"`
		Deaths   []item `json:"deaths"`
		Holidays []item `json:"holidays"`
	}

	if err := json.Unmarshal(body, &feed); err != nil {
		return "", "", "", nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	// collect pages together with their originating item's evidence
	type entry struct {
		p page
		e Evidence
	}
	entries := []entry{}
	addItems := func(items []item) {
		for _, it := range items {
			ev := Evidence{Type: it.Type, Year: it.Year, Text: it.Text}
			for _, pg := range it.Pages {
				entries = append(entries, entry{p: pg, e: ev})
			}
		}
	}
	addItems(feed.Events)
	addItems(feed.Births)
	addItems(feed.Deaths)
	addItems(feed.Holidays)

	if len(entries) == 0 {
		return "", "", "", nil, fmt.Errorf("no pages found for date")
	}

	// pick random
	idx := int(time.Now().UnixNano() % int64(len(entries)))
	ent := entries[idx]
	return ent.p.Title, ent.p.Extract, ent.p.ContentURLs.Desktop.Page, &ent.e, nil
}

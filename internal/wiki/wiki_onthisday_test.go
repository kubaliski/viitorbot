package wiki

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetRandomArticleByDate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// minimal feed with one event containing one page
		_, _ = w.Write([]byte(`{
            "events": [ { "pages": [ { "title": "Evento", "extract": "Extracto evento", "content_urls": { "desktop": { "page": "https://es.wikipedia.org/wiki/Evento" } } } ] } ],
            "births": [], "deaths": [], "holidays": []
        }`))
	}))
	defer ts.Close()

	old := OnThisDayEndpoint
	OnThisDayEndpoint = ts.URL + "/%02d/%02d"
	defer func() { OnThisDayEndpoint = old }()

	c := NewClient()
	title, extract, url, evidence, err := c.GetRandomArticleByDate(time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "Evento" {
		t.Fatalf("expected title Evento, got %q", title)
	}
	if !strings.Contains(url, "Evento") {
		t.Fatalf("unexpected url: %s", url)
	}
	if extract == "" {
		t.Fatalf("expected extract")
	}
	if evidence == nil {
		t.Fatalf("expected evidence for onthisday entry")
	}
}

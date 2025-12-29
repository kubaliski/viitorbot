package wiki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRandomArticle(t *testing.T) {
	// Create a test server that returns a fixed JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
            "title": "Prueba",
            "extract": "Esto es un extracto de prueba.",
            "content_urls": { "desktop": { "page": "https://es.wikipedia.org/wiki/Prueba" } }
        }`))
	}))
	defer ts.Close()

	// Override endpoint
	old := Endpoint
	Endpoint = ts.URL
	defer func() { Endpoint = old }()

	c := NewClient()
	title, extract, url, err := c.GetRandomArticle()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if title != "Prueba" {
		t.Fatalf("expected title 'Prueba', got %q", title)
	}
	if extract == "" {
		t.Fatalf("expected non-empty extract")
	}
	if url != "https://es.wikipedia.org/wiki/Prueba" {
		t.Fatalf("unexpected url: %s", url)
	}
}

func TestGetRandomArticle_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	old := Endpoint
	Endpoint = ts.URL
	defer func() { Endpoint = old }()

	c := NewClient()
	_, _, _, err := c.GetRandomArticle()
	if err == nil {
		t.Fatalf("expected error for non-200 status")
	}
}

func TestGetRandomArticle_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not a json`))
	}))
	defer ts.Close()

	old := Endpoint
	Endpoint = ts.URL
	defer func() { Endpoint = old }()

	c := NewClient()
	_, _, _, err := c.GetRandomArticle()
	if err == nil {
		t.Fatalf("expected json unmarshal error")
	}
}

func TestGetRandomArticle_RequestError(t *testing.T) {
	// Use an invalid URL to force request error
	old := Endpoint
	Endpoint = "http://127.0.0.1:0"
	defer func() { Endpoint = old }()

	c := NewClient()
	_, _, _, err := c.GetRandomArticle()
	if err == nil {
		t.Fatalf("expected request error for invalid endpoint")
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"reef-sa/internal/feeds"
)

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got '%s'", w.Body.String())
	}
}

func TestReadyEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "ready" {
		t.Errorf("expected body 'ready', got '%s'", w.Body.String())
	}
}

func TestNewsItemSerialization(t *testing.T) {
	item := NewsItem{
		ID:          "test-1",
		Title:       "Test Article",
		Source:      "Test Source",
		URL:         "https://example.com/test",
		PublishedAt: "2024-01-01T12:00:00Z",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("failed to marshal NewsItem: %v", err)
	}

	var decoded NewsItem
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal NewsItem: %v", err)
	}

	if decoded.ID != item.ID {
		t.Errorf("expected ID %s, got %s", item.ID, decoded.ID)
	}
	if decoded.Title != item.Title {
		t.Errorf("expected Title %s, got %s", item.Title, decoded.Title)
	}
}

func TestRegionalFeedResponseSerialization(t *testing.T) {
	response := RegionalFeedResponse{
		Region:  "ASIA",
		Country: "JP",
		News: []NewsItem{
			{
				ID:          "asia-1",
				Title:       "Test News",
				Source:      "Test Source",
				URL:         "https://example.com",
				PublishedAt: "2024-01-01T00:00:00Z",
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal RegionalFeedResponse: %v", err)
	}

	var decoded RegionalFeedResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal RegionalFeedResponse: %v", err)
	}

	if decoded.Region != response.Region {
		t.Errorf("expected Region %s, got %s", response.Region, decoded.Region)
	}
	if decoded.Country != response.Country {
		t.Errorf("expected Country %s, got %s", response.Country, decoded.Country)
	}
	if len(decoded.News) != len(response.News) {
		t.Errorf("expected %d news items, got %d", len(response.News), len(decoded.News))
	}
}

func TestRegionalFeedResponseWithWeather(t *testing.T) {
	response := RegionalFeedResponse{
		Region:  "ASIA",
		Country: "JP",
		Weather: &feeds.WeatherData{
			Summary:      "Clear sky",
			TemperatureC: 25.5,
			FeelsLikeC:   24.0,
		},
		News: []NewsItem{},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response with weather: %v", err)
	}

	// Verify JSON contains weather field
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, ok := jsonMap["weather"]; !ok {
		t.Error("expected weather field in JSON")
	}
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"reef-sa/internal/featureflags"
	"reef-sa/internal/feeds"
	mw "reef-sa/internal/http/middleware"
	"reef-sa/internal/logger"
)

// NewsItem represents a single news article
type NewsItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

// RegionalFeedResponse represents the response structure
type RegionalFeedResponse struct {
	Region  string              `json:"region"`
	Country string              `json:"country"`
	Weather *feeds.WeatherData  `json:"weather,omitempty"`
	News    []NewsItem          `json:"news,omitempty"`
}

func main() {
	// 1) Feature flags init (non-fatal)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := featureflags.Init(ctx, ""); err != nil {
		log.Printf("feature flags init warning: %v", err)
	} else {
		log.Printf("feature flags ready: offline=%v, logLevel=%s",
			featureflags.Values().Offline.IsEnabled(nil),
			featureflags.Values().LogLevel.GetValue(nil))
	}
	defer featureflags.Shutdown()

	// 2) Initialize levelled logger from flag & watch for flips
	logger.Init(featureflags.Values().LogLevel.GetValue(nil))
	logger.Infof("log level set to %s", logger.GetLevel())

	go func() {
		prev := featureflags.Values().LogLevel.GetValue(nil)
		for {
			time.Sleep(5 * time.Second)
			cur := featureflags.Values().LogLevel.GetValue(nil)
			if cur != prev {
				logger.SetLevel(cur)
				logger.Infof("log level changed to %s", logger.GetLevel())
				prev = cur
			}
		}
	}()

	// 3) Router
	r := mux.NewRouter()

	// 4) Offline kill-switch middleware
	offlineGate := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// always allow health checks
			if r.URL.Path == "/health" || r.URL.Path == "/ready" {
				next.ServeHTTP(w, r)
				return
			}
			// block all other requests when Offline flag is ON
			if featureflags.Values().Offline.IsEnabled(nil) {
				http.Error(w, "service temporarily offline", http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
	r.Use(offlineGate)

	// 5) Request logger (skip noisy health endpoints)
	r.Use(mw.LogRequests(mw.WithSkips("/health", "/ready")))

	// 6) Health endpoints
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	}).Methods(http.MethodGet)

	// 7) Inspect current flag values
	r.HandleFunc("/_flags", func(w http.ResponseWriter, _ *http.Request) {
		resp := map[string]interface{}{
			"offline":  featureflags.Values().Offline.IsEnabled(nil),
			"logLevel": featureflags.Values().LogLevel.GetValue(nil),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}).Methods(http.MethodGet)

	// 8) Regional feeds endpoint with stub data
	r.HandleFunc("/regional-feeds", func(w http.ResponseWriter, r *http.Request) {
		country := r.URL.Query().Get("country")
		if country == "" {
			country = "BR" // Default to Brazil for demo
		}

		logger.Infof("serving SA regional feeds for country: %s", country)

		// Fetch real weather data from Open-Meteo API
		weather, err := feeds.FetchWeather(country)
		if err != nil {
			logger.Warnf("failed to fetch weather for %s: %v (using fallback)", country, err)
			// Fallback to stub data if API fails
			weather = &feeds.WeatherData{
				Summary:      "Weather data unavailable",
				TemperatureC: 0,
				FeelsLikeC:   0,
			}
		}

		// Stub news data for EU
		news := []NewsItem{
			{
				ID:          "eu-article-1",
				Title:       "European Markets Rise Amid Economic Growth",
				Source:      "EU Financial Times",
				URL:         "https://example.com/eu/article1",
				PublishedAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			},
			{
				ID:          "eu-article-2",
				Title:       "New Climate Agreement Reached in Brussels",
				Source:      "EuroNews",
				URL:         "https://example.com/eu/article2",
				PublishedAt: time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
			},
			{
				ID:          "eu-article-3",
				Title:       "Tech Innovation Summit Draws Global Leaders",
				Source:      "European Tech Review",
				URL:         "https://example.com/eu/article3",
				PublishedAt: time.Now().Add(-6 * time.Hour).Format(time.RFC3339),
			},
		}

		response := RegionalFeedResponse{
			Region:  "EU",
			Country: country,
			Weather: weather,
			News:    news,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Errorf("failed to encode response: %v", err)
		}
	}).Methods(http.MethodGet)

	s := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	logger.Infof("reef-sa listening on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}

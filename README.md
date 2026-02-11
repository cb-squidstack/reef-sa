# Reef EU - Europe Regional Service

Regional feed service for Europe in the SquidStack Reef Feeds feature. Provides weather and news data for European countries.

## Overview

The reef-sa service provides:
- Stub weather data for European countries
- Stub news data with European headlines
- Simple HTTP JSON API

## Endpoints

- `GET /health` - Liveness probe
- `GET /ready` - Readiness probe
- `GET /_flags` - Inspect feature flag values
- `GET /regional-feeds?country=XX` - Get regional feeds for country (e.g., ?country=GB)

## Configuration

### Environment Variables

- `FM_NAMESPACE` - CloudBees Feature Management namespace (default: "default")

### Feature Management Key

Mount FM key file at `/app/config/fm.json`

## Feature Flags

- `offline` - Kill switch to put service offline
- `logLevel` - Log level (debug, info, warn, error)

## Development

```bash
# Build
go build -o app main.go

# Run
./app

# Docker build
docker build -t reef-sa .

# Docker run
docker run -p 8080:8080 reef-sa
```

## API Response Example

```json
{
  "region": "EU",
  "country": "GB",
  "weather": {
    "summary": "Cloudy",
    "temperatureC": 14.2,
    "feelsLikeC": 12.5
  },
  "news": [
    {
      "id": "eu-article-1",
      "title": "European Markets Rise Amid Economic Growth",
      "source": "EU Financial Times",
      "url": "https://example.com/eu/article1",
      "publishedAt": "2025-12-11T10:30:00Z"
    }
  ]
}
```

Last updated: 2026-02-11 21:51 UTC

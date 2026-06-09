package weatherapi

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// NewProviderFromEnv selects the weather data provider.
// WEATHER_PROVIDER: "stub" (default) or "http".
func NewProviderFromEnv() (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("WEATHER_PROVIDER"))) {
	case "", "stub":
		return StubProvider{}, nil
	case "http":
		return &HTTPProvider{
			BaseURL: os.Getenv("WEATHER_API_BASE_URL"),
			APIKey:  os.Getenv("WEATHER_API_KEY"),
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown WEATHER_PROVIDER: %s", os.Getenv("WEATHER_PROVIDER"))
	}
}

package weatherapi

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ErrExternalProviderNotConfigured is returned until a real upstream API is wired in.
var ErrExternalProviderNotConfigured = errors.New("external weather provider not configured")

// HTTPProvider is a placeholder for a future real-world weather API client.
type HTTPProvider struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func (p *HTTPProvider) Name() string { return "http" }

// Lookup will call an external weather API once configured.
// Set WEATHER_API_BASE_URL and WEATHER_API_KEY when enabling this provider.
func (p *HTTPProvider) Lookup(ctx context.Context, tz string, at time.Time) (Forecast, error) {
	_ = ctx
	_ = tz
	_ = at
	_ = p.BaseURL
	_ = p.APIKey
	return Forecast{}, ErrExternalProviderNotConfigured
}

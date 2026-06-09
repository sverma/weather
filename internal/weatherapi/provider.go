package weatherapi

import (
	"context"
	"time"
)

// Conditions describes high-level weather state.
type Conditions struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Temperature holds values in both unit systems.
type Temperature struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
}

// Wind describes wind at observation time.
type Wind struct {
	SpeedKPH  float64 `json:"speed_kph"`
	Direction string  `json:"direction"`
}

// Forecast is the weather payload returned by providers.
type Forecast struct {
	Timezone      string      `json:"timezone"`
	LocalDatetime string      `json:"local_datetime"`
	Location      string      `json:"location"`
	Conditions    Conditions  `json:"conditions"`
	Temperature   Temperature `json:"temperature"`
	HumidityPct   int         `json:"humidity_percent"`
	Wind          Wind        `json:"wind"`
	Provider      string      `json:"provider"`
	ObservedAt    string      `json:"observed_at"`
}

// Provider fetches weather for an IANA timezone at a reference instant.
type Provider interface {
	Name() string
	Lookup(ctx context.Context, tz string, at time.Time) (Forecast, error)
}

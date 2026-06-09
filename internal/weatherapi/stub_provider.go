package weatherapi

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

var stubConditionCatalog = []Conditions{
	{Summary: "Sunny", Description: "Clear skies and bright sunshine", Icon: "sunny"},
	{Summary: "Partly cloudy", Description: "A mix of sun and clouds", Icon: "partly_cloudy"},
	{Summary: "Cloudy", Description: "Overcast with grey skies", Icon: "cloudy"},
	{Summary: "Light rain", Description: "Intermittent light showers", Icon: "light_rain"},
	{Summary: "Rainy", Description: "Steady rainfall expected", Icon: "rainy"},
	{Summary: "Windy", Description: "Breezy conditions throughout the day", Icon: "windy"},
}

var stubWindDirections = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

// StubProvider returns deterministic placeholder weather derived from the timezone name.
type StubProvider struct{}

func (StubProvider) Name() string { return "stub" }

func (StubProvider) Lookup(_ context.Context, tz string, at time.Time) (Forecast, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return Forecast{}, fmt.Errorf("invalid timezone: %s", tz)
	}

	seed := hashTimezone(tz)
	condition := stubConditionCatalog[seed%uint32(len(stubConditionCatalog))]
	tempC := 5.0 + float64(seed%300)/10.0
	humidity := 30 + int(seed%61)
	windKPH := 3.0 + float64(seed%280)/10.0
	windDir := stubWindDirections[seed%uint32(len(stubWindDirections))]

	local := at.In(loc)
	return Forecast{
		Timezone:      tz,
		LocalDatetime: local.Format(time.RFC3339),
		Location:      locationFromTimezone(tz),
		Conditions:    condition,
		Temperature:   Temperature{Celsius: tempC, Fahrenheit: celsiusToFahrenheit(tempC)},
		HumidityPct:   humidity,
		Wind: Wind{
			SpeedKPH:  windKPH,
			Direction: windDir,
		},
		Provider:   "stub",
		ObservedAt: at.UTC().Format(time.RFC3339),
	}, nil
}

func hashTimezone(tz string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(tz))
	return h.Sum32()
}

func locationFromTimezone(tz string) string {
	parts := strings.Split(tz, "/")
	last := parts[len(parts)-1]
	return strings.ReplaceAll(last, "_", " ")
}

func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

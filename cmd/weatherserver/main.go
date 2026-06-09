package main

import (
	"log"
	"net/http"
	"os"

	"weather/internal/observability"
	"weather/internal/weatherapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	provider, err := weatherapi.NewProviderFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	metrics := observability.New(observability.BuildInfoFromEnv())
	h := weatherapi.NewHandler(provider).WithMetrics(metrics)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /metrics", metrics.Instrument(observability.RouteMetrics, func(w http.ResponseWriter, r *http.Request) {
		metrics.MetricsHandler().ServeHTTP(w, r)
	}))
	mux.HandleFunc("GET /health", metrics.Instrument(observability.RouteHealth, observability.HealthHandler))
	mux.HandleFunc("GET /ready", metrics.Instrument(observability.RouteReady, observability.ReadyHandler))
	mux.HandleFunc("GET /weather", metrics.Instrument(observability.RouteWeather, h.Weather))

	addr := ":" + port
	log.Printf("weather service listening on %s (provider=%s)", addr, provider.Name())
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

package observability

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsHandlerExposesREDMetrics(t *testing.T) {
	c := New(BuildInfo{Version: "test", GoVersion: "go1.24", GitCommit: "abc"})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /metrics", c.Instrument(RouteMetrics, func(w http.ResponseWriter, r *http.Request) {
		c.MetricsHandler().ServeHTTP(w, r)
	}))
	mux.HandleFunc("GET /ok", c.Instrument("/ok", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	mux.HandleFunc("GET /bad", c.Instrument("/bad", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/ok", nil))
	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/bad", nil))

	c.RecordTimezoneLookupError("invalid_tz")
	c.RecordProviderRequest("stub")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body, _ := io.ReadAll(rec.Body)
	output := string(body)

	for _, want := range []string{
		"http_requests_total",
		"http_request_errors_total",
		"http_request_duration_seconds",
		"weather_timezone_lookup_errors_total",
		"weather_provider_requests_total",
		"weather_build_info",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("metrics output missing %q", want)
		}
	}
}

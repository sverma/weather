package observability

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// BuildInfo identifies the running binary for GitOps verification.
type BuildInfo struct {
	Version   string
	GoVersion string
	GitCommit string
}

// BuildInfoFromEnv reads VERSION and GIT_COMMIT from the environment.
func BuildInfoFromEnv() BuildInfo {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "dev"
	}
	gitCommit := os.Getenv("GIT_COMMIT")
	if gitCommit == "" {
		gitCommit = "unknown"
	}
	return BuildInfo{
		Version:   version,
		GoVersion: runtime.Version(),
		GitCommit: gitCommit,
	}
}

// Collector owns Prometheus metrics and HTTP instrumentation helpers.
type Collector struct {
	registry *prometheus.Registry

	requestsTotal      *prometheus.CounterVec
	requestErrorsTotal *prometheus.CounterVec
	requestDuration    *prometheus.HistogramVec
	requestsInFlight   prometheus.Gauge

	timezoneLookupErrors *prometheus.CounterVec
	providerRequests     *prometheus.CounterVec
	providerErrors       *prometheus.CounterVec
	buildInfo            *prometheus.GaugeVec
	processStartTime     prometheus.Gauge
}

// New creates a Collector and registers application and runtime metrics.
func New(info BuildInfo) *Collector {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	c := &Collector{
		registry: reg,

		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests.",
		}, []string{"method", "route", "status_class"}),

		requestErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total HTTP request errors.",
		}, []string{"method", "route", "error_type"}),

		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		}, []string{"method", "route"}),

		requestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being served.",
		}),

		timezoneLookupErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_timezone_lookup_errors_total",
			Help: "Timezone validation failures on weather endpoints.",
		}, []string{"reason"}),

		providerRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_provider_requests_total",
			Help: "Successful weather provider lookups.",
		}, []string{"provider"}),

		providerErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_provider_errors_total",
			Help: "Failed weather provider lookups.",
		}, []string{"provider", "reason"}),

		buildInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "weather_build_info",
			Help: "Build and version information for the running binary.",
		}, []string{"version", "go_version", "git_commit"}),

		processStartTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "weather_process_start_time_seconds",
			Help: "Unix timestamp when the process started.",
		}),
	}

	reg.MustRegister(
		c.requestsTotal,
		c.requestErrorsTotal,
		c.requestDuration,
		c.requestsInFlight,
		c.timezoneLookupErrors,
		c.providerRequests,
		c.providerErrors,
		c.buildInfo,
		c.processStartTime,
	)

	c.buildInfo.WithLabelValues(info.Version, info.GoVersion, info.GitCommit).Set(1)
	c.processStartTime.Set(float64(time.Now().Unix()))

	return c
}

// MetricsHandler serves Prometheus metrics at GET /metrics.
func (c *Collector) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
}

func (c *Collector) RecordTimezoneLookupError(reason string) {
	c.timezoneLookupErrors.WithLabelValues(reason).Inc()
}

func (c *Collector) RecordProviderRequest(provider string) {
	c.providerRequests.WithLabelValues(provider).Inc()
}

func (c *Collector) RecordProviderError(provider, reason string) {
	c.providerErrors.WithLabelValues(provider, reason).Inc()
}

package weatherapi

// MetricsRecorder records application-level metrics for weather handlers.
type MetricsRecorder interface {
	RecordTimezoneLookupError(reason string)
	RecordProviderRequest(provider string)
	RecordProviderError(provider, reason string)
}

type noopMetricsRecorder struct{}

func (noopMetricsRecorder) RecordTimezoneLookupError(string)   {}
func (noopMetricsRecorder) RecordProviderRequest(string)       {}
func (noopMetricsRecorder) RecordProviderError(string, string) {}

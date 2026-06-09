package weatherapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Handler serves weather HTTP endpoints.
type Handler struct {
	clock    Clock
	provider Provider
	metrics  MetricsRecorder
}

// NewHandler returns a handler with the system clock and given provider.
func NewHandler(provider Provider) *Handler {
	return &Handler{
		clock:    RealClock{},
		provider: provider,
		metrics:  noopMetricsRecorder{},
	}
}

// NewHandlerWithClock returns a handler with a custom clock (for tests).
func NewHandlerWithClock(provider Provider, clock Clock) *Handler {
	return &Handler{
		clock:    clock,
		provider: provider,
		metrics:  noopMetricsRecorder{},
	}
}

// WithMetrics attaches a metrics recorder to the handler.
func (h *Handler) WithMetrics(recorder MetricsRecorder) *Handler {
	h.metrics = recorder
	return h
}

type errorResponse struct {
	Error string `json:"error"`
}

// Weather handles GET /weather — current weather for an IANA timezone.
//
// Query parameters:
//   - tz: IANA timezone (required), e.g. Europe/London
//   - at: optional reference instant (RFC3339); defaults to now
func (h *Handler) Weather(w http.ResponseWriter, r *http.Request) {
	tz := strings.TrimSpace(r.URL.Query().Get("tz"))
	if tz == "" {
		h.recordTimezoneError("missing_param")
		writeError(w, http.StatusBadRequest, "tz is required")
		return
	}

	at, errMsg, status := parseReferenceInstant(r.URL.Query().Get("at"), h.clock)
	if errMsg != "" {
		writeError(w, status, errMsg)
		return
	}

	forecast, err := h.provider.Lookup(r.Context(), tz, at)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid timezone") {
			h.recordTimezoneError("invalid_tz")
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.recordProviderError(h.provider.Name(), "lookup_failed")
		writeError(w, http.StatusServiceUnavailable, "weather provider unavailable")
		return
	}

	h.recordProviderRequest(h.provider.Name())
	writeJSON(w, http.StatusOK, forecast)
}

func parseReferenceInstant(at string, clock Clock) (time.Time, string, int) {
	at = strings.TrimSpace(at)
	if at == "" {
		return clock.Now().UTC(), "", http.StatusOK
	}

	instant, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return time.Time{}, "invalid at: " + at, http.StatusBadRequest
	}
	return instant.UTC(), "", http.StatusOK
}

func (h *Handler) recordTimezoneError(reason string) {
	if h.metrics != nil {
		h.metrics.RecordTimezoneLookupError(reason)
	}
}

func (h *Handler) recordProviderRequest(provider string) {
	if h.metrics != nil {
		h.metrics.RecordProviderRequest(provider)
	}
}

func (h *Handler) recordProviderError(provider, reason string) {
	if h.metrics != nil {
		h.metrics.RecordProviderError(provider, reason)
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

package observability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	HealthHandler(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	var body probeResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" {
		t.Errorf("status = %q", body.Status)
	}
}

func TestReadyHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	ReadyHandler(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))

	var body probeResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ready" {
		t.Errorf("status = %q", body.Status)
	}
}

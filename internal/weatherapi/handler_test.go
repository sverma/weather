package weatherapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (c fixedClock) Now() time.Time {
	return c.t
}

var testInstant = time.Date(2026, 6, 3, 14, 30, 45, 0, time.UTC)

func newTestHandler() *Handler {
	return NewHandlerWithClock(StubProvider{}, fixedClock{t: testInstant})
}

func TestWeatherSuccess(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/weather?tz=Europe/London", nil)
	rec := httptest.NewRecorder()

	h.Weather(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body Forecast
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Timezone != "Europe/London" {
		t.Errorf("timezone = %q", body.Timezone)
	}
	if body.Location != "London" {
		t.Errorf("location = %q, want London", body.Location)
	}
	if body.LocalDatetime != "2026-06-03T15:30:45+01:00" {
		t.Errorf("local_datetime = %q", body.LocalDatetime)
	}
	if body.Provider != "stub" {
		t.Errorf("provider = %q, want stub", body.Provider)
	}
}

func TestWeatherMissingTZ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	h.Weather(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestWeatherInvalidTZ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/weather?tz=Bad/Zone", nil)
	rec := httptest.NewRecorder()

	h.Weather(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestWeatherWithAt(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/weather?tz=UTC&at=2026-01-15T12:00:00Z", nil)
	rec := httptest.NewRecorder()

	h.Weather(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestStubProviderDeterministic(t *testing.T) {
	p := StubProvider{}
	a, err := p.Lookup(nil, "Asia/Tokyo", testInstant)
	if err != nil {
		t.Fatal(err)
	}
	b, err := p.Lookup(nil, "Asia/Tokyo", testInstant)
	if err != nil {
		t.Fatal(err)
	}
	if a.Conditions.Summary != b.Conditions.Summary {
		t.Fatal("expected deterministic stub data for same timezone")
	}
}

# weather — timezone weather sample service

A standalone HTTP microservice that returns weather details for an IANA timezone. Designed to be called by **gtod** (and other services) over HTTP. Uses **stub data** today; swap to a real external API later without changing the HTTP contract.

## Service layout

```
/Users/saurabhverma/code/
├── gtod/       # time service (port 8080)
└── weather/    # this service (port 8081)
```

## Primary API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/weather` | Current weather for an IANA timezone |

### Query parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `tz` | yes | IANA timezone, e.g. `Europe/London`, `Asia/Tokyo` |
| `at` | no | Reference instant (RFC3339); defaults to now |

### Example

```bash
curl -s "http://localhost:8081/weather?tz=Europe/London"
```

```json
{
  "timezone": "Europe/London",
  "local_datetime": "2026-06-03T15:30:45+01:00",
  "location": "London",
  "conditions": {
    "summary": "Partly cloudy",
    "description": "A mix of sun and clouds",
    "icon": "partly_cloudy"
  },
  "temperature": {
    "celsius": 18.5,
    "fahrenheit": 65.3
  },
  "humidity_percent": 62,
  "wind": {
    "speed_kph": 12.4,
    "direction": "SW"
  },
  "provider": "stub",
  "observed_at": "2026-06-03T14:30:45Z"
}
```

Stub data is **deterministic per timezone** (same `tz` always returns the same conditions until you switch providers).

## Observability

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/metrics` | Prometheus metrics (RED + runtime + build info) |
| `GET` | `/health` | Liveness — `{"status":"ok"}` |
| `GET` | `/ready` | Readiness — `{"status":"ready"}` |

### RED metrics

| Metric | Type | Labels |
|--------|------|--------|
| `http_requests_total` | Counter | `method`, `route`, `status_class` |
| `http_request_errors_total` | Counter | `method`, `route`, `error_type` |
| `http_request_duration_seconds` | Histogram | `method`, `route` |
| `http_requests_in_flight` | Gauge | — |

### Application metrics

| Metric | Type | Labels |
|--------|------|--------|
| `weather_timezone_lookup_errors_total` | Counter | `reason` |
| `weather_provider_requests_total` | Counter | `provider` |
| `weather_provider_errors_total` | Counter | `provider`, `reason` |
| `weather_build_info` | Gauge | `version`, `go_version`, `git_commit` |
| `weather_process_start_time_seconds` | Gauge | — |

## Provider configuration

| Env var | Default | Description |
|---------|---------|-------------|
| `WEATHER_PROVIDER` | `stub` | `stub` or `http` |
| `WEATHER_API_BASE_URL` | — | Future external API base URL |
| `WEATHER_API_KEY` | — | Future external API key |
| `PORT` | `8081` | HTTP listen port |
| `VERSION` | `dev` | Build version for `weather_build_info` |
| `GIT_COMMIT` | `unknown` | Git commit for `weather_build_info` |

`http` provider is a **placeholder** — it returns `503` until a real upstream client is implemented in `internal/weatherapi/http_provider.go`.

## Integration with gtod

gtod calls weather over HTTP (service-to-service):

```text
Client → gtod:8080/time?tz=Europe/London
gtod   → weather:8081/weather?tz=Europe/London   (future gtod endpoint or internal client)
```

Suggested gtod env var:

```bash
WEATHER_SERVICE_URL=http://weather:8081
```

gtod can proxy or enrich responses, e.g. `GET /time/weather?tz=Europe/London` combining local time + weather.

## Build & run

```bash
go test ./... -cover -count=1
go build -o bin/weatherserver ./cmd/weatherserver
./bin/weatherserver

# or
PORT=8081 go run ./cmd/weatherserver
```

## Project layout

```
.
├── cmd/weatherserver/       # HTTP server entrypoint
├── internal/weatherapi/     # Handlers, stub/http providers
├── internal/observability/  # Prometheus metrics and probes
├── Dockerfile
├── go.mod
└── README.md
```

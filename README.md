# onboard

An HR-Portal.

## Server

HTTP API built with:

- [gorilla/mux](https://github.com/gorilla/mux) — routing
- [pgx](https://github.com/jackc/pgx) — PostgreSQL
- [zap](https://github.com/uber-go/zap) — structured logging
- [viper](https://github.com/spf13/viper) — configuration (YAML + env)

### Configuration

Defaults live in `configs/config.yaml`. Environment variables override the file.
Copy `.env.example` to `.env` for local overrides (loaded by Viper).

| Variable | Description |
| --- | --- |
| `SERVER_HOST` | Bind address (default `0.0.0.0`) |
| `SERVER_PORT` | Bind port (default `8080`) |
| `DATABASE_URL` | PostgreSQL DSN |
| `LOG_LEVEL` | `debug`, `info`, `warn`, or `error` |
| `LOG_ENCODING` | `json` or `console` |
| `LOG_DEVELOPMENT` | `true` for development logger |

### Run

PostgreSQL must be reachable at `DATABASE_URL` before the process starts.

```bash
go run ./cmd/server
```

### Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Liveness |
| `GET` | `/ready` | Readiness (pings PostgreSQL) |
| `GET` | `/api/v1/` | Service info |

### Tests

```bash
go test ./...
```

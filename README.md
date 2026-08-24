# onboard

An HR-Portal.

## Server

HTTP API built with:

- [gorilla/mux](https://github.com/gorilla/mux) — routing
- [pgx](https://github.com/jackc/pgx) — PostgreSQL
- [zap](https://github.com/uber-go/zap) — structured JSON logging
- [viper](https://github.com/spf13/viper) — YAML configuration

### Configuration

All settings come from `configs/config.yaml`. There are no defaults and no environment overlays — edit the YAML file before running.

### Run

PostgreSQL must be reachable at `database.url` in the config file.

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

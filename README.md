# onboard

An HR-Portal.

## Server

HTTP API built with:

- [gorilla/mux](https://github.com/gorilla/mux) — routing
- [pgx](https://github.com/jackc/pgx) — PostgreSQL
- [zap](https://github.com/uber-go/zap) — structured JSON logging
- [viper](https://github.com/spf13/viper) — YAML configuration

### Configuration

Copy the sample file and edit local values. `configs/config.yaml` is gitignored.

```bash
cp configs/config.sample.yaml configs/config.yaml
```

Viper reads only `configs/config.yaml`. There are no defaults and no environment overlays.

### Run

PostgreSQL must be reachable at `database.url` in the local config file.

```bash
cp configs/config.sample.yaml configs/config.yaml
go run ./cmd/server
```

### Tests

```bash
go test ./...
```

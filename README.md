# driftwatch

Lightweight daemon that detects config drift between deployed services and source definitions.

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build ./...
```

## Usage

Point driftwatch at your config source and a running service endpoint:

```bash
driftwatch --source ./configs/service.yaml --target http://localhost:8080/config --interval 30s
```

driftwatch will poll on the defined interval and report any fields that have drifted from the source definition.

### Example output

```
[2024-01-15 09:42:01] DRIFT DETECTED in service: api-gateway
  - max_connections: expected=500, got=250
  - timeout_ms: expected=3000, got=5000
[2024-01-15 09:42:31] OK: no drift detected
```

### Configuration file

```yaml
source: ./configs/
target: http://localhost:8080/config
interval: 30s
alert_webhook: https://hooks.example.com/notify
```

Run as a daemon:

```bash
driftwatch --config driftwatch.yaml --daemon
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--source` | `./configs` | Path to source config definitions |
| `--target` | — | Service endpoint to check against |
| `--interval` | `60s` | Poll interval |
| `--daemon` | `false` | Run as background daemon |

## License

MIT
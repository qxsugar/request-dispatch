# Request Dispatch

A powerful Traefik middleware plugin for intelligent request routing based on request markers, enabling gray-scale deployments, controlled feature rollouts, and system upgrades.

## Features

- **Intelligent Request Routing** — Route requests to different backend services based on HTTP headers
- **Gray-Scale Deployments** — Gradually roll out new versions to a subset of users
- **Controlled Feature Rollouts** — Test new features with specific user groups before full deployment
- **Load Balancing** — Randomly distribute requests among multiple backend hosts for a single mark
- **Easy Configuration** — Simple YAML-based configuration with Traefik

## Installation

Add the plugin to your Traefik configuration:

```yaml
experimental:
  plugins:
    request-dispatch:
      moduleName: github.com/qxsugar/request-dispatch
      version: v1.0.2
```

## Configuration

### Basic Setup

```yaml
http:
  routers:
    api:
      rule: host(`api.example.com`)
      service: production-service
      entryPoints:
        - web
      middlewares:
        - request-dispatch

  services:
    production-service:
      loadBalancer:
        servers:
          - url: "https://prod.api.example.com"

  middlewares:
    request-dispatch:
      plugin:
        request-dispatch:
          logLevel: INFO
          markHeader: X-Dispatch-Mark
          markHosts:
            canary:
              - https://canary-v2.api.example.com
              - https://canary-v2-backup.api.example.com
            staging:
              - https://staging.api.example.com
```

### Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `logLevel` | string | Log level: `DEBUG`, `INFO`, or `ERROR` (case-insensitive) |
| `markHeader` | string | HTTP header name to check for routing marks |
| `markHosts` | map | Mapping of mark values to lists of backend URLs |

## Usage

### Basic Request Routing

Send requests with the mark header to route them to specific backends:

```bash
# Route to canary backend
curl -H "X-Dispatch-Mark: canary" https://api.example.com/users

# Route to staging backend
curl -H "X-Dispatch-Mark: staging" https://api.example.com/users

# Route to production (default)
curl https://api.example.com/users
```

### Load Balancing

When multiple hosts are configured for a mark, the plugin randomly selects one:

```yaml
markHosts:
  canary:
    - https://canary-v2.api.example.com
    - https://canary-v2-backup.api.example.com  # Requests randomly distributed
```

### Gray-Scale Deployment Example

```yaml
markHosts:
  v2:
    - https://api-v2.example.com
  v1:
    - https://api-v1.example.com
```

Route 10% of users to v2 for testing:
```bash
# For 10% of requests
curl -H "X-Dispatch-Mark: v2" https://api.example.com/...

# For remaining 90% (default route)
curl https://api.example.com/...
```

## Development

### Build

```bash
go build
```

### Test

```bash
make test          # Run all tests with coverage
make test -v       # Run tests with verbose output
go test -v -run TestDispatch ./...  # Run specific test
```

### Lint

```bash
make lint
```

### Dependencies

```bash
make vendor        # Vendor dependencies
make clean         # Remove vendor directory
```

## Architecture

The plugin intercepts HTTP requests and:

1. Checks for the configured mark header
2. If the header matches a configured mark, randomly selects a backend host
3. Proxies the request to the selected backend
4. Falls back to the default route if no mark is found or header is missing

### Thread Safety

The random host selection is protected by a mutex to ensure thread-safe operation under concurrent load.

### Error Handling

If reverse proxy fails, the request automatically falls back to the default handler, ensuring graceful degradation.

## Logging

The plugin supports three log levels:

- **DEBUG** — Detailed information about request routing and proxy operations
- **INFO** — General information about plugin operations
- **ERROR** — Error messages only

Log level is case-insensitive in configuration.

## License

See LICENSE file for details.

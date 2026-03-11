# goforge

The fastest way to start an idiomatic production Go service.

goforge scaffolds team-ready Go HTTP services with strong defaults, low magic, and fast time-to-first-run. Generated code compiles and runs immediately — no cleanup needed.

## Install

```bash
go install github.com/byzkhan/goforge/cmd/goforge@latest
```

## Usage

### Generate a new service

```bash
# Interactive (recommended)
goforge new service myapi

# Non-interactive with defaults
goforge new service myapi --non-interactive

# Fully specified
goforge new service myapi \
  --module github.com/mycompany/myapi \
  --router chi \
  --db postgres \
  --observability otel \
  --auth jwt-stub \
  --ci github \
  --docker
```

### Other commands

```bash
goforge version    # Print version
goforge doctor     # Check your environment
```

## Flags

| Flag | Options | Default |
|------|---------|---------|
| `--module` | Go module path | `github.com/yourorg/<name>` |
| `--router` | `chi`, `stdlib` | `chi` |
| `--db` | `none`, `postgres` | `none` |
| `--observability` | `none`, `otel` | `none` |
| `--auth` | `none`, `jwt-stub` | `none` |
| `--ci` | `github`, `none` | `github` |
| `--docker` | `true`, `false` | `true` |
| `--non-interactive` | Skip prompts | `false` |

## What you get

A complete, runnable Go HTTP service:

```
myapi/
├── cmd/myapi/main.go           # Entrypoint
├── internal/
│   ├── app/app.go              # Dependency wiring
│   ├── config/config.go        # Env-based config
│   ├── http/
│   │   ├── handlers/           # Request handlers
│   │   ├── middleware/         # Request ID, logging, panic recovery
│   │   ├── routes/            # Route registration
│   │   └── response.go       # JSON response helpers
│   ├── domain/                # Domain types
│   ├── service/               # Business logic
│   ├── repository/            # Data access
│   └── platform/
│       ├── logger/            # Structured logging (slog)
│       ├── server/            # HTTP server with timeouts
│       └── shutdown/          # Graceful shutdown
├── tests/                     # Integration tests
├── Makefile                   # run/build/test/lint
├── Dockerfile                 # Multi-stage build
├── .github/workflows/ci.yml  # GitHub Actions
└── go.mod
```

### Included out of the box

- Health (`/healthz`) and readiness (`/readyz`) endpoints
- Graceful shutdown on SIGTERM/SIGINT
- Structured JSON logging with `slog`
- Request ID middleware
- Request logging middleware
- Panic recovery middleware
- Production HTTP server timeouts
- Consistent JSON error responses
- Example CRUD endpoints
- Unit and integration tests
- Makefile with `run`, `build`, `test`, `lint`
- Multi-stage Dockerfile
- GitHub Actions CI

### Optional features

- **PostgreSQL** (`--db postgres`): pgx connection pool, migrations, repository pattern
- **OpenTelemetry** (`--observability otel`): tracing with stdout exporter (swap for OTLP in production)
- **JWT auth** (`--auth jwt-stub`): HMAC-SHA256 Bearer token middleware (replace with your auth provider)

## Design principles

- **Low magic**: No reflection, no hidden registries, no DI containers
- **Explicit wiring**: Dependencies are constructed and passed via constructors in `app.go`
- **Minimal dependencies**: Only what's justified (chi for routing, pgx for postgres, otel for tracing)
- **Production-shaped**: The generated code has the same structure you'd use at scale
- **Easy to extend**: Clear separation of concerns without over-engineered abstractions
- **Fast to start**: `cd myapi && make run` — that's it

## Development

```bash
make build    # Build goforge
make test     # Run tests
make install  # Install to GOPATH/bin
```

## License

MIT

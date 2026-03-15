# Gateway

Schema-first HTTP gateway built on Go 1.22 with OpenAPI-driven request and response validation.

Architecture: [Gateway](https://github.com/agynio/architecture/blob/main/architecture/gateway.md)

## Local Development

Full setup: [Local Development](https://github.com/agynio/architecture/blob/main/architecture/operations/local-development.md)

### Prepare environment

```bash
git clone https://github.com/agynio/bootstrap.git
cd bootstrap
chmod +x apply.sh
./apply.sh -y
```

See [bootstrap](https://github.com/agynio/bootstrap) for details.

### Run from sources

```bash
# Deploy once (exit when healthy)
devspace dev

# Watch mode (streams logs, re-syncs on changes)
devspace dev -w
```

### Run tests

```bash
devspace run test:e2e
```

See [E2E Testing](https://github.com/agynio/architecture/blob/main/architecture/operations/e2e-testing.md).

## Adding a New API Domain

Every API domain in the gateway must be spec-driven and implemented through
oapi-codegen strict servers (no hand-written CRUD handlers). The standard flow
mirrors the Team API:

1. Pull the OpenAPI artifact into `.openapi/` via `scripts/pull-spec.sh`, then
   extend `scripts/generate-openapi.sh` to embed it under
   `internal/apischema/<domain>v1/` and generate strict server stubs into
   `internal/<domain>gen/`, followed by `gofmt`.
2. Implement `internal/handlers/<domain>.go` to satisfy the generated
   `StrictServerInterface`, with conversions in
   `internal/handlers/<domain>_convert.go` and direct handler tests in
   `internal/handlers/<domain>_test.go`.
3. Wire the router in `cmd/gateway/main.go` with request validation middleware
   and optional response validation for CRUD endpoints. Streaming/proxy routes
   (like SSE `/responses`) must still be defined in the spec for request
   validation but mounted as raw handlers outside the strict server group.
4. Update CI to generate OpenAPI code and build Redoc docs for the new spec.

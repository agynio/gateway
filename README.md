# Gateway API

Schema-first HTTP gateway built on Go 1.22 with OpenAPI-driven request and response validation.

## Specification

The Team and LLM OpenAPI 3.0 definitions are pulled from `ghcr.io/agynio/openapi/{team,llm}:1` into `.openapi/` and generated at build time. Use `scripts/generate-openapi.sh` to pull the latest specs, embed them under `internal/apischema/`, and regenerate the `oapi-codegen` server stubs under `internal/gen` and `internal/llmgen`.

## Prerequisites

- Go 1.22+
- [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) and [`oras`](https://github.com/oras-project/oras) on your `$PATH` for local OpenAPI generation
- [`schemathesis`](https://github.com/schemathesis/schemathesis) and [`redocly`](https://github.com/Redocly/redoc) for validation/docs (installed automatically in CI)

## Quickstart

```bash
TEAMS_GRPC_TARGET="localhost:50051" \
FILES_GRPC_TARGET="localhost:50052" \
LLM_GRPC_TARGET="localhost:50053" \
SECRETS_GRPC_TARGET="localhost:50054" \
go run ./cmd/gateway
```

In Kubernetes, the gRPC targets default to the `teams`, `files`, `llm`, and `secrets` services on port 50051. Override them locally as needed. The Team Management API remains available under `/team/v1` for backwards compatibility.

Enable optional response validation middleware by exporting `OPENAPI_VALIDATE_RESPONSE=true` before starting the server.

## Container image

Pre-built multi-architecture images are published to GitHub Container Registry as [`ghcr.io/agynio/gateway`](https://ghcr.io/agynio/gateway).

- Pushes to `main` publish the `main` and `sha-<shortSHA>` tags.
- Git tags following `vX.Y.Z` publish the `X.Y.Z`, `X.Y`, `X`, and `latest` tags.

Run the gateway container locally by supplying the required gRPC targets:

```bash
docker run --rm -p 8080:8080 \
  -e TEAMS_GRPC_TARGET="localhost:50051" \
  -e FILES_GRPC_TARGET="localhost:50052" \
  -e LLM_GRPC_TARGET="localhost:50053" \
  -e SECRETS_GRPC_TARGET="localhost:50054" \
  ghcr.io/agynio/gateway:main
```

## Helm chart

An OCI-packaged Helm chart is available at `oci://ghcr.io/agynio/charts/gateway`. The chart delegates to the `service-base` library and exposes typed values for all gateway configuration knobs.

Install the chart by wiring any non-default gRPC targets and response validation flags as needed:

```bash
helm install gateway oci://ghcr.io/agynio/charts/gateway \
  --version 0.1.0 \
  --namespace gateway --create-namespace \
  --set gateway.teamsGrpcTarget="teams:50051" \
  --set gateway.filesGrpcTarget="files:50051"
```

Additional gRPC targets and response validation flags can be adjusted via the `gateway.*` values or by overriding `env`/`extraEnvVars` directly.

## Development

Regenerate the embedded specs and Chi servers whenever the OpenAPI documents change:

```bash
bash scripts/generate-openapi.sh
```

Regenerate gRPC stubs as needed:

```bash
buf generate buf.build/agynio/api \
  --path agynio/api/files/v1 \
  --path agynio/api/llm/v1 \
  --path agynio/api/secrets/v1 \
  --path agynio/api/teams/v1
```

Run the full local test suite (mirrors CI):

```bash
buf generate buf.build/agynio/api \
  --path agynio/api/files/v1 \
  --path agynio/api/llm/v1 \
  --path agynio/api/secrets/v1 \
  --path agynio/api/teams/v1
bash scripts/generate-openapi.sh
go vet ./...
go test ./...
redocly build-docs .openapi/team-v1.yaml --output dist/redoc.html
redocly build-docs .openapi/llm-v1.yaml --output dist/llm-redoc.html
schemathesis run --url=http://localhost:8080 .openapi/team-v1.yaml
```

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

## Continuous Integration

The CI workflow performs:

1. `buf generate` for gRPC stubs.
2. OpenAPI generation (`scripts/generate-openapi.sh`).
3. `go vet` and `go test` over the module.
4. Redocly build with the rendered documentation uploaded as an artifact.

Artifacts and reports are available on each run via the GitHub Actions UI.

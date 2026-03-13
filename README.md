# Gateway API

Schema-first HTTP gateway built on Go 1.22 with OpenAPI-driven request and response validation.

## Specification

The Team API OpenAPI 3.0 definition is bundled as [`internal/apischema/teamv1/team-v1.yaml`](internal/apischema/teamv1/team-v1.yaml) and pulled from `ghcr.io/agynio/openapi/team:1`. Use `scripts/pull-spec.sh` to fetch the latest spec into `.openapi/team-v1.yaml`, then run `scripts/sync-embedded-spec.sh` to update the embedded copy. Generated server stubs and models are located under [`internal/gen`](internal/gen).

## Prerequisites

- Go 1.22+
- [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) binary on your `$PATH`
- [`spectral`](https://github.com/stoplightio/spectral), [`oasdiff`](https://github.com/Tufin/oasdiff), [`schemathesis`](https://github.com/schemathesis/schemathesis), and [`redocly`](https://github.com/Redocly/redoc) for local validation (installed automatically in CI)

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

Regenerate the typed models and Chi server whenever the OpenAPI document changes:

```bash
bash scripts/pull-spec.sh
oapi-codegen --config oapi-codegen.server.yaml .openapi/team-v1.yaml
gofmt -w internal/gen/server.gen.go
bash scripts/sync-embedded-spec.sh
```

Run the full local test suite (mirrors CI):

```bash
spectral lint .openapi/team-v1.yaml
oasdiff breaking --fail-on ERR .openapi/team-v1.yaml .openapi/team-v1.yaml
go test ./...
go vet ./...
schemathesis run --url=http://localhost:8080 .openapi/team-v1.yaml
```

## Adding a New API Domain

Every API domain in the gateway must be spec-driven and implemented through
oapi-codegen strict servers (no hand-written CRUD handlers). The standard flow
mirrors the Team API:

1. Pull the OpenAPI artifact into `.openapi/` via `scripts/pull-spec.sh`, then
   sync it into `internal/apischema/<domain>v1/` with
   `scripts/sync-embedded-spec.sh`.
2. Add an `oapi-codegen.<domain>.yaml` config and generate strict server stubs
   into `internal/<domain>gen/`, followed by `gofmt`.
3. Implement `internal/handlers/<domain>.go` to satisfy the generated
   `StrictServerInterface`, with conversions in
   `internal/handlers/<domain>_convert.go` and direct handler tests in
   `internal/handlers/<domain>_test.go`.
4. Wire the router in `cmd/gateway/main.go` with request validation middleware
   and optional response validation for CRUD endpoints. Streaming/proxy routes
   (like SSE `/responses`) must still be defined in the spec for request
   validation but mounted as raw handlers outside the strict server group.
5. Update CI to lint, diff, re-generate, and build Redoc docs for the new spec.

## Continuous Integration

The CI workflow performs:

1. Spectral linting of the OpenAPI description.
2. `oasdiff` breaking-change guard against the main branch specification.
3. Re-generation checks for `oapi-codegen` output and Go format drift.
4. `go vet` and `go test` over the module.
5. Redocly build with the rendered documentation uploaded as an artifact.

Artifacts and reports are available on each run via the GitHub Actions UI.

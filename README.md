# Gateway API

Schema-first HTTP gateway built on Go 1.22 with OpenAPI-driven request and response validation.

## Specification

The OpenAPI 3.0 definition lives in [`spec/openapi.yaml`](spec/openapi.yaml). Generated server stubs and models are located under [`internal/gen`](internal/gen).

## Prerequisites

- Go 1.22+
- [`oapi-codegen`](https://github.com/oapi-codegen/oapi-codegen) binary on your `$PATH`
- [`spectral`](https://github.com/stoplightio/spectral), [`oasdiff`](https://github.com/Tufin/oasdiff), [`schemathesis`](https://github.com/schemathesis/schemathesis), and [`redocly`](https://github.com/Redocly/redoc) for local validation (installed automatically in CI)

## Quickstart

```bash
PLATFORM_BASE_URL="https://platform.example.com" \
PLATFORM_AUTH_TOKEN="<token>" \
go run ./cmd/gateway
```

All Team Management endpoints are served under the `/team/v1` prefix and proxy to the configured platform server.

Enable optional response validation middleware by exporting `OPENAPI_VALIDATE_RESPONSE=true` before starting the server.

## Development

Regenerate the typed models and Chi server whenever the OpenAPI document changes:

```bash
oapi-codegen --config oapi-codegen.server.yaml spec/openapi.yaml
gofmt -w internal/gen/server.gen.go
```

Run the full local test suite (mirrors CI):

```bash
spectral lint spec/openapi.yaml
oasdiff breaking --fail-on ERR spec/openapi.yaml spec/openapi.yaml
go test ./...
go vet ./...
schemathesis run --url=http://localhost:8080 spec/openapi.yaml
```

## Continuous Integration

The CI workflow performs:

1. Spectral linting of the OpenAPI description.
2. `oasdiff` breaking-change guard against the main branch specification.
3. Re-generation checks for `oapi-codegen` output and Go format drift.
4. `go vet` and `go test` over the module.
5. Redocly build with the rendered documentation uploaded as an artifact.

Artifacts and reports are available on each run via the GitHub Actions UI.

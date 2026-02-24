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
go run ./cmd/gateway
```

### Health check

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### Personalized greeting

```bash
curl \
  -H 'Content-Type: application/json' \
  -d '{"name":"Casey"}' \
  http://localhost:8080/hello
# {"message":"Hello, Casey"}
```

Enable optional response validation middleware by exporting `VALIDATE_RESPONSES=true` before starting the server.

## Development

Regenerate the typed models and Chi server whenever the OpenAPI document changes:

```bash
oapi-codegen --package=gen --generate types -o internal/gen/types.gen.go spec/openapi.yaml
oapi-codegen --config oapi-codegen.server.yaml spec/openapi.yaml
go fmt ./internal/gen/...
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
5. Schemathesis contract tests against the running service.
6. Redocly build with the rendered documentation uploaded as an artifact.

Artifacts and reports are available on each run via the GitHub Actions UI.

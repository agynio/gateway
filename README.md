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
PLATFORM_BASE_URL="https://api.agyn.dev:8080" \
PLATFORM_AUTH_TOKEN="<token>" \
go run ./cmd/gateway
```

The gateway forwards `/health` and the `/api/*` surface directly to the configured platform server. The existing Team Management API remains available under `/team/v1` for backwards compatibility.

Enable optional response validation middleware by exporting `OPENAPI_VALIDATE_RESPONSE=true` before starting the server.

## Container image

Pre-built multi-architecture images are published to GitHub Container Registry as [`ghcr.io/agynio/gateway`](https://ghcr.io/agynio/gateway).

- Pushes to `main` publish the `main` and `sha-<shortSHA>` tags.
- Git tags following `vX.Y.Z` publish the `X.Y.Z`, `X.Y`, `X`, and `latest` tags.

Run the gateway container locally by supplying the required platform settings:

```bash
docker run --rm -p 8080:8080 \
  -e PLATFORM_BASE_URL="https://api.agyn.dev:8080" \
  -e PLATFORM_AUTH_TOKEN="<token>" \
  ghcr.io/agynio/gateway:main
```

## Helm chart

An OCI-packaged Helm chart is available at `oci://ghcr.io/agynio/charts/gateway`. The chart delegates to the `service-base` library and exposes typed values for all gateway configuration knobs.

Install the chart by providing the target platform URL (required) and optionally wiring a secret for the auth token:

```bash
helm install gateway oci://ghcr.io/agynio/charts/gateway \
  --version 0.1.0 \
  --namespace gateway --create-namespace \
  --set gateway.platformBaseUrl="https://api.agyn.dev:8080" \
  --set gateway.authToken.existingSecret=gateway-auth \
  --set gateway.authToken.existingSecretKey=token
```

Additional timeouts, retry counts, request headers, and response validation flags can be adjusted via the `gateway.*` values or by overriding `env`/`extraEnvVars` directly.

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

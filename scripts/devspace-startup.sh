#!/usr/bin/env bash
set -eu

echo "=== DevSpace startup ==="

export PATH="$(go env GOPATH)/bin:$PATH"
ORAS_VERSION="1.1.0"
OAPI_CODEGEN_VERSION="v2.4.0"

if ! command -v oras >/dev/null 2>&1; then
  echo "Installing oras..."
  go install "oras.land/oras/cmd/oras@v${ORAS_VERSION}"
fi

if ! command -v oapi-codegen >/dev/null 2>&1; then
  echo "Installing oapi-codegen..."
  go install "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}"
fi

echo "Generating protobuf types..."
buf generate buf.build/agynio/api --path agynio/api/files/v1 --path agynio/api/llm/v1 --path agynio/api/secrets/v1 --path agynio/api/teams/v1

echo "Generating OpenAPI code..."
bash scripts/generate-openapi.sh

echo "Downloading Go modules..."
go mod download

echo "Starting dev server (air)..."
exec air

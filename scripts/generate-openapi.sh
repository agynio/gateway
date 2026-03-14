#!/usr/bin/env bash
# generate-openapi.sh - Pull OpenAPI specs from GHCR and generate Go server code.
# Requires: oras, oapi-codegen, gofmt on $PATH.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# -- Step 1: Pull specs ------------------------------------------------
bash "${SCRIPT_DIR}/pull-spec.sh"

# -- Step 2: Place specs for //go:embed -------------------------------
mkdir -p "${ROOT_DIR}/internal/apischema/teamv1"
mkdir -p "${ROOT_DIR}/internal/apischema/llmv1"
cp "${ROOT_DIR}/.openapi/team-v1.yaml" "${ROOT_DIR}/internal/apischema/teamv1/team-v1.yaml"
cp "${ROOT_DIR}/.openapi/llm-v1.yaml"  "${ROOT_DIR}/internal/apischema/llmv1/llm-v1.yaml"

# -- Step 3: Generate spec.go loaders ---------------------------------
cat > "${ROOT_DIR}/internal/apischema/teamv1/spec.go" << 'GOEOF'
package teamv1

import (
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed team-v1.yaml
var Spec []byte

func LoadSpec() (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	return loader.LoadFromData(Spec)
}
GOEOF

cat > "${ROOT_DIR}/internal/apischema/llmv1/spec.go" << 'GOEOF'
package llmv1

import (
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed llm-v1.yaml
var Spec []byte

func LoadSpec() (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	return loader.LoadFromData(Spec)
}
GOEOF

# -- Step 4: Run oapi-codegen -----------------------------------------
mkdir -p "${ROOT_DIR}/internal/llmgen"

oapi-codegen \
  -package gen \
  -generate chi-server,strict-server,models \
  -o "${ROOT_DIR}/internal/gen/server.gen.go" \
  "${ROOT_DIR}/.openapi/team-v1.yaml"

oapi-codegen \
  -package llmgen \
  -generate chi-server,strict-server,models \
  -o "${ROOT_DIR}/internal/llmgen/server.gen.go" \
  "${ROOT_DIR}/.openapi/llm-v1.yaml"

# -- Step 5: Format generated code ------------------------------------
gofmt -w "${ROOT_DIR}/internal/gen/server.gen.go"
gofmt -w "${ROOT_DIR}/internal/llmgen/server.gen.go"

echo "OpenAPI generation complete."

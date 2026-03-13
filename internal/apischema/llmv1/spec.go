package llmv1

import (
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
)

// Spec holds the bundled LLM API v1 OpenAPI specification (single YAML file,
// all $refs resolved). Pulled from ghcr.io/agynio/openapi/llm:1.
//
//go:embed llm-v1.yaml
var Spec []byte

// LoadSpec parses the embedded LLM API v1 OpenAPI spec.
func LoadSpec() (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	return loader.LoadFromData(Spec)
}

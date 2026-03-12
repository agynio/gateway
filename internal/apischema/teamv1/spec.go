package teamv1

import (
	_ "embed"

	"github.com/getkin/kin-openapi/openapi3"
)

// Spec holds the bundled Team API v1 OpenAPI specification (single YAML file,
// all $refs resolved). Pulled from ghcr.io/agynio/openapi/team:1.
//
//go:embed team-v1.yaml
var Spec []byte

// LoadSpec parses the embedded Team API v1 OpenAPI spec.
func LoadSpec() (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	return loader.LoadFromData(Spec)
}

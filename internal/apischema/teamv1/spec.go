package teamv1

import _ "embed"

// Spec holds the bundled Team API v1 OpenAPI specification (single YAML file,
// all $refs resolved). Pulled from ghcr.io/agynio/openapi/team:1.
//
//go:embed team-v1.yaml
var Spec []byte

package teamv1

import "embed"

// SpecFS bundles the canonical Team API specification copied from
// github.com/agynio/api during go:generate. Handlers rely on this
// read-only filesystem to load schemas for runtime validation.
//
//go:embed openapi.yaml components/**/*.yaml paths/*.yaml
var SpecFS embed.FS

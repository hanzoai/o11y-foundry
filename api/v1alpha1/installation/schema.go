package installation

import (
	_ "embed"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
)

//go:embed casting.schema.json
var schemaJSON []byte

var schema = v1alpha1.MustResolveSchema(schemaJSON)

// Schema returns the resolved JSON schema for an Installation casting.
func Schema() *jsonschema.Resolved {
	return schema
}

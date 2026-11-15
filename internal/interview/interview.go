// Package interview validates the report an interview agent writes for one Advisor.
package interview

import (
	"bytes"
	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// schemaURL names the schema inside the compiler and in error messages.
const schemaURL = "urn:agneta:interview"

//go:embed schema.json
var schemaJSON []byte

var schema = compileSchema()

// Schema returns a copy of the embedded schema, for publishing it unchanged.
func Schema() []byte {
	return bytes.Clone(schemaJSON)
}

// compileSchema compiles the embedded schema with format assertion on, and panics if it is invalid.
func compileSchema() *jsonschema.Schema {
	parsed, err := parseJSON(schemaJSON)
	if err != nil {
		panic(err)
	}
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	if err := compiler.AddResource(schemaURL, parsed); err != nil {
		panic(err)
	}
	return compiler.MustCompile(schemaURL)
}

// Validate returns nil for a valid report, else the parse error or a *jsonschema.ValidationError naming every bad field.
func Validate(doc []byte) error {
	parsed, err := parseJSON(doc)
	if err != nil {
		return err
	}
	return schema.Validate(parsed)
}

// parseJSON decodes raw JSON into the form the jsonschema package validates.
func parseJSON(raw []byte) (any, error) {
	return jsonschema.UnmarshalJSON(bytes.NewReader(raw))
}

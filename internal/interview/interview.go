// Package interview checks the document an interview agent writes after
// interviewing one person.
package interview

import (
	"bytes"
	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schema.json
var schemaJSON []byte

var schema = mustCompile(schemaJSON)

// mustCompile compiles a JSON Schema with format assertion on, and panics if
// the schema itself is invalid.
func mustCompile(source []byte) *jsonschema.Schema {
	parsed, err := jsonschema.UnmarshalJSON(bytes.NewReader(source))
	if err != nil {
		panic(err)
	}
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat()
	const location = "urn:agneta:interview"
	if err := compiler.AddResource(location, parsed); err != nil {
		panic(err)
	}
	return compiler.MustCompile(location)
}

// Validate reports whether doc, a JSON document, is a valid interview document.
// It returns nil for a valid document and an error describing every problem otherwise.
func Validate(doc []byte) error {
	parsed, err := jsonschema.UnmarshalJSON(bytes.NewReader(doc))
	if err != nil {
		return err
	}
	return schema.Validate(parsed)
}

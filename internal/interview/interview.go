// Package interview checks the document an interview agent writes after
// interviewing one person.
package interview

import (
	"bytes"
	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// schemaURL names the interview schema inside the compiler. It shows up in
// validation error messages.
const schemaURL = "urn:agneta:interview"

//go:embed schema.json
var schemaJSON []byte

var schema = compileSchema()

// compileSchema compiles the embedded interview schema with format assertion
// on. It panics if the schema itself is invalid.
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

// Validate reports whether doc, a JSON document, is a valid interview document.
// It returns nil for a valid document. Otherwise it returns the JSON parse
// error, or a *jsonschema.ValidationError listing every problem with the path
// to each bad field.
func Validate(doc []byte) error {
	parsed, err := parseJSON(doc)
	if err != nil {
		return err
	}
	return schema.Validate(parsed)
}

// parseJSON decodes raw JSON into the value form the jsonschema package
// validates. It returns the parse error if raw is not valid JSON.
func parseJSON(raw []byte) (any, error) {
	return jsonschema.UnmarshalJSON(bytes.NewReader(raw))
}

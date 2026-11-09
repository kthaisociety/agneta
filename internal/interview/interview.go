// Package interview checks the document an interview agent writes after
// interviewing one person.
package interview

import "errors"

// Validate reports whether doc, a JSON document, is a valid interview document.
// It returns nil for a valid document and an error describing every problem otherwise.
func Validate(doc []byte) error {
	return errors.New("not implemented")
}

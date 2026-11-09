package interview

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// validDoc returns a fresh interview document that the schema should accept,
// for a test to break in one place.
func validDoc() map[string]any {
	return map[string]any{
		"name": "Ada Lovelace",
		"date": "2026-11-09",
		"role": "Board advisor",
		"segments": []any{
			map[string]any{
				"question": "Who signs off on the yearly budget?",
				"answer":   "The board votes on it in the January meeting.",
			},
		},
	}
}

// firstSegment returns the first Q&A segment of doc, for a test to change.
func firstSegment(doc map[string]any) map[string]any {
	return doc["segments"].([]any)[0].(map[string]any)
}

// validate encodes doc as JSON and runs it through Validate.
func validate(t *testing.T, doc map[string]any) error {
	t.Helper()
	encoded, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("encode document: %v", err)
	}
	return Validate(encoded)
}

// failure is one problem the schema found: where in the document, and which
// schema keyword it broke.
type failure struct {
	location string
	keyword  string
}

// failures lists every problem a schema validation error holds.
func failures(err *jsonschema.ValidationError) []failure {
	if len(err.Causes) == 0 {
		location := ""
		if len(err.InstanceLocation) > 0 {
			location = "/" + strings.Join(err.InstanceLocation, "/")
		}
		keywordPath := err.ErrorKind.KeywordPath()
		keyword := ""
		if len(keywordPath) > 0 {
			keyword = keywordPath[len(keywordPath)-1]
		}
		return []failure{{location, keyword}}
	}
	var all []failure
	for _, cause := range err.Causes {
		all = append(all, failures(cause)...)
	}
	return all
}

// assertRejected fails the test unless doc is rejected by the schema with a
// problem at location that breaks keyword.
func assertRejected(t *testing.T, doc map[string]any, location, keyword string) {
	t.Helper()
	err := validate(t, doc)
	var validationErr *jsonschema.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("want a schema validation error, got %v", err)
	}
	want := failure{location, keyword}
	got := failures(validationErr)
	for _, f := range got {
		if f == want {
			return
		}
	}
	t.Fatalf("want a %q failure at %q, got %+v", keyword, location, got)
}

// TestValidate_AcceptsDocument checks that the schema accepts a document with
// name, date and role at the top and every question paired with its answer.
func TestValidate_AcceptsDocument(t *testing.T) {
	doc := validDoc()
	doc["role"] = "Tech lead"
	doc["segments"] = append(doc["segments"].([]any), map[string]any{
		"question": "Who owns the sponsor contacts?",
		"answer":   "The head of partnerships.",
	})
	if err := validate(t, doc); err != nil {
		t.Fatalf("want the document accepted, got %v", err)
	}
}

// TestValidate_AcceptsExample checks that the three-Q&A example document
// passes the schema.
func TestValidate_AcceptsExample(t *testing.T) {
	example, err := os.ReadFile("example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	var parsed struct {
		Segments []any `json:"segments"`
	}
	if err := json.Unmarshal(example, &parsed); err != nil {
		t.Fatalf("parse example: %v", err)
	}
	if len(parsed.Segments) != 3 {
		t.Fatalf("want 3 Q&A segments in the example, got %d", len(parsed.Segments))
	}
	if err := Validate(example); err != nil {
		t.Fatalf("want the example accepted, got %v", err)
	}
}

// TestValidate_RejectsEmptySegmentList checks that a document with no Q&A
// segments is rejected.
func TestValidate_RejectsEmptySegmentList(t *testing.T) {
	doc := validDoc()
	doc["segments"] = []any{}
	assertRejected(t, doc, "/segments", "minItems")
}

// TestValidate_RejectsMissingTopLevelField checks that a document missing its
// name, date or role is rejected.
func TestValidate_RejectsMissingTopLevelField(t *testing.T) {
	for _, field := range []string{"name", "date", "role"} {
		t.Run(field, func(t *testing.T) {
			doc := validDoc()
			delete(doc, field)
			assertRejected(t, doc, "", "required")
		})
	}
}

// TestValidate_RejectsBlankName checks that an empty or whitespace-only name
// is rejected.
func TestValidate_RejectsBlankName(t *testing.T) {
	for _, name := range []string{"", " ", " \t\n "} {
		t.Run(name, func(t *testing.T) {
			doc := validDoc()
			doc["name"] = name
			assertRejected(t, doc, "/name", "pattern")
		})
	}
}

// TestValidate_RejectsImpossibleDate checks that a date that is not a real
// calendar date is rejected.
func TestValidate_RejectsImpossibleDate(t *testing.T) {
	for _, date := range []string{"2026-02-30", "2026-13-01", "9 November 2026", "yesterday"} {
		t.Run(date, func(t *testing.T) {
			doc := validDoc()
			doc["date"] = date
			assertRejected(t, doc, "/date", "format")
		})
	}
}

// TestValidate_RejectsUnknownRole checks that a role outside Board advisor,
// Board member and Tech lead is rejected.
func TestValidate_RejectsUnknownRole(t *testing.T) {
	for _, role := range []string{"CEO", "board advisor", ""} {
		t.Run(role, func(t *testing.T) {
			doc := validDoc()
			doc["role"] = role
			assertRejected(t, doc, "/role", "enum")
		})
	}
}

// TestValidate_RejectsSegmentMissingField checks that a segment missing its
// question or answer is rejected.
func TestValidate_RejectsSegmentMissingField(t *testing.T) {
	for _, field := range []string{"question", "answer"} {
		t.Run(field, func(t *testing.T) {
			doc := validDoc()
			delete(firstSegment(doc), field)
			assertRejected(t, doc, "/segments/0", "required")
		})
	}
}

// TestValidate_RejectsBlankSegmentField checks that an empty or
// whitespace-only question or answer is rejected.
func TestValidate_RejectsBlankSegmentField(t *testing.T) {
	for _, field := range []string{"question", "answer"} {
		for _, value := range []string{"", " ", " \t\n "} {
			t.Run(field+"/"+value, func(t *testing.T) {
				doc := validDoc()
				firstSegment(doc)[field] = value
				assertRejected(t, doc, "/segments/0/"+field, "pattern")
			})
		}
	}
}

// TestValidate_RejectsUnknownField checks that an unknown extra field, at the
// top level or in a segment, is rejected.
func TestValidate_RejectsUnknownField(t *testing.T) {
	t.Run("top level", func(t *testing.T) {
		doc := validDoc()
		doc["topic"] = "budget"
		assertRejected(t, doc, "", "additionalProperties")
	})
	t.Run("segment", func(t *testing.T) {
		doc := validDoc()
		firstSegment(doc)["topic"] = "budget"
		assertRejected(t, doc, "/segments/0", "additionalProperties")
	})
}

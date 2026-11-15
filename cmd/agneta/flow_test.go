package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kthaisociety/agneta/internal/interview"
)

// connect serves the service's handler and opens an MCP session to it, as a connector would.
func connect(t *testing.T) *mcp.ClientSession {
	t.Helper()
	service := httptest.NewServer(handler(slog.New(slog.NewTextHandler(io.Discard, nil))))
	t.Cleanup(service.Close)
	client := mcp.NewClient(&mcp.Implementation{Name: "flow-test", Version: "0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.StreamableClientTransport{Endpoint: service.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

// submit calls submit_report with doc and returns the result's text and whether it is an error.
func submit(t *testing.T, session *mcp.ClientSession, doc []byte) (string, bool) {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: "submit_report", Arguments: json.RawMessage(doc)})
	if err != nil {
		t.Fatalf("call submit_report: %v", err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("result has %d content items, want 1", len(result.Content))
	}
	return result.Content[0].(*mcp.TextContent).Text, result.IsError
}

// TestFlowListsSubmitReportWithTheEmbeddedSchema checks that the tool's input schema is schema.json unchanged.
func TestFlowListsSubmitReportWithTheEmbeddedSchema(t *testing.T) {
	tools, err := connect(t).ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "submit_report" {
		t.Fatalf("tools = %v, want only submit_report", tools.Tools)
	}
	if got, want := canonical(t, tools.Tools[0].InputSchema), canonical(t, json.RawMessage(interview.Schema())); !bytes.Equal(got, want) {
		t.Errorf("input schema differs from schema.json:\n got %s\nwant %s", got, want)
	}
}

// TestFlowAcceptsAValidReport checks that the example report is accepted, counted and marked not stored.
func TestFlowAcceptsAValidReport(t *testing.T) {
	doc, err := os.ReadFile("../../internal/interview/example.json")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	text, isError := submit(t, connect(t), doc)
	if isError || text != `{"accepted":true,"stored":false,"record_count":3}` {
		t.Errorf("result = %s (error %v)", text, isError)
	}
}

// TestFlowNamesEveryProblemWithAnInvalidReport checks that an invalid report is a tool error naming each bad field.
func TestFlowNamesEveryProblemWithAnInvalidReport(t *testing.T) {
	text, isError := submit(t, connect(t), []byte(`{"name": " ", "date": "2026-02-30", "role": "CEO",
		"segments": [{"question": "q", "answer": ""}], "author": "x"}`))
	if !isError {
		t.Fatalf("invalid report accepted: %s", text)
	}
	for _, path := range []string{"'/name'", "'/date'", "'/role'", "'/segments/0/answer'", "'author'"} {
		if !strings.Contains(text, path) {
			t.Errorf("error does not name %s:\n%s", path, text)
		}
	}
}

// canonical re-encodes a JSON value so two schemas compare regardless of formatting.
func canonical(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("decode: %v", err)
	}
	out, _ := json.Marshal(parsed)
	return out
}

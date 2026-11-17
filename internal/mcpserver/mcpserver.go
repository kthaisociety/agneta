// Package mcpserver serves agneta's MCP tools over streamable HTTP (docs/adr/0001).
package mcpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kthaisociety/agneta/internal/interview"
)

// SubmitReport is the tool the interview skill calls with an approved report.
const SubmitReport = "submit_report"

const submitDescription = "Submit one interview report to agneta, only after the Advisor has approved every " +
	"question and answer. The arguments are the report itself. If the result lists problems, fix only what it " +
	"names, never invent content, and call it again. agneta doesn't store reports yet: a valid report comes " +
	"back with stored: false."

// receipt is the result of a valid submission.
type receipt struct {
	Accepted    bool `json:"accepted"`
	Stored      bool `json:"stored"`
	RecordCount int  `json:"record_count"`
}

// Handler returns the stateless MCP endpoint with the submit_report tool.
func Handler(logger *slog.Logger) http.Handler {
	server := mcp.NewServer(&mcp.Implementation{Name: "agneta", Version: "0.1.0"}, nil)
	// The low-level AddTool publishes schema.json unchanged and leaves validation to interview.Validate.
	server.AddTool(&mcp.Tool{
		Name:        SubmitReport,
		Description: submitDescription,
		InputSchema: json.RawMessage(interview.Schema()),
	}, submitReport(logger))
	return mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true},
	)
}

// submitReport validates the arguments and acknowledges a valid report; problems come back as tool errors.
func submitReport(logger *slog.Logger) mcp.ToolHandler {
	return func(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		doc := []byte(req.Params.Arguments)
		if err := interview.Validate(doc); err != nil {
			return failed("The report is invalid. Fix every problem below and call submit_report again.\n" + err.Error()), nil
		}
		var report struct {
			Segments []json.RawMessage `json:"segments"`
		}
		if err := json.Unmarshal(doc, &report); err != nil {
			return nil, err
		}
		logger.Info("report accepted", "records", len(report.Segments))
		body, err := json.Marshal(receipt{Accepted: true, Stored: false, RecordCount: len(report.Segments)})
		if err != nil {
			return nil, err
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(body)}}}, nil
	}
}

// failed returns a tool result the agent sees as a failed call.
func failed(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}

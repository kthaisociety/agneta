# 0001. Submit reports through an MCP tool

Status: Accepted

## Context
The interview skill (PR #5) ends with an approved report that the Advisor was told to upload somewhere. The
agent running the interview already holds that report, and MCP is how an agent calls a service's tools. If
the agent submits it, the Advisor copies nothing, and the agent can read validation errors and fix them with
the Advisor.

## Decision
- agneta is an MCP server: no separate upload service. It serves `/mcp` over streamable HTTP with the official
  Go SDK (PCR-0001), stateless, so every request stands alone.
- One tool, `submit_report`. Its input schema is `internal/interview/schema.json` exactly as embedded,
  registered with the SDK's low-level `Server.AddTool` so the SDK neither derives a schema nor validates. The
  handler validates with `interview.Validate` (PCR-0003, PCR-0004).
- An invalid report is a tool result with `isError` listing every bad field's path, not a protocol error, so
  the agent reads it.
- For now there is no sign-in and no storage, so the skill-to-connector flow can be exercised end to end
  before either is designed. A valid report is answered with `{"accepted": true, "stored": false,
  "record_count": N}` and nothing is kept.

## Consequences
- The skill's last step becomes "call `submit_report`".
- Anyone who can reach `/mcp` can call it. Until sign-in exists, it only runs locally or behind a short-lived
  tunnel, with synthetic interviews.
- Sign-in through identity-service, storage in Postgres, the author, and resubmitting without storing twice
  each get their own decision when built.

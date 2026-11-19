# agneta
Capture what advisors and the board know but never wrote down, so they can be offboarded and the board still has guidance.

---

## Development

The backend is Go 1.27 or later, an MCP server with one tool, `submit_report` (docs/adr/0001). Run the tests the way the PR check does with `make test`.

Run it locally:

```
cp .env.example .env
make run
```

Start a local Postgres and create its tables (the settings come from `.env`):

```
docker compose up -d
go tool goose up
```

The server doesn't use the database yet; the tables are described in docs/adr/0002.

The MCP server serves `GET /health`, which answers `ok`, and the MCP endpoint at `/mcp`. There is no sign-in and no storage connected yet: a valid report is answered with `{"accepted": true, "stored": false, "record_count": N}`, and an invalid one with a tool error listing every bad field. Use synthetic interviews only.

### Connect Claude

- **Claude Code** reaches localhost directly: `claude mcp add --transport http agneta http://127.0.0.1:8003/mcp`.
- **claude.ai, including voice in the mobile app**, calls connectors from Anthropic's servers, so it needs a public URL. Open a tunnel with `cloudflared tunnel --url http://127.0.0.1:8003 --http-host-header 127.0.0.1:8003` (the MCP SDK answers 403 to a tunnel's own `Host` header, to block DNS rebinding), then add a custom connector at claude.ai with the tunnel's URL plus `/mcp`. It then shows up in the mobile app too. Anyone with the URL can call the tool, so close the tunnel when you're done.

The interview schema lives in `internal/interview/schema.json`, with an example report next to it in `example.json`.

# AGENTS.md

agneta: captures what KTHAIS Advisors know but never wrote down, so the Board keeps that guidance after they
leave. The MVP only collects knowledge: an interview skill writes a JSON report, the upload service
validates it and stores its records in Postgres. Nothing reads them back yet.

## Commands

- Test: `go test ./...`
- Run: `go run ./cmd/agneta` (port from `PORT`, default 8000)
- Image: `docker build -t agneta .`

## MVP

- **Interview skill**: an agent skill, outside this repo, that runs an interview and writes one report.
- **Report**: one JSON file per interview: the interviewee's name, role and the date once at the top, then a
  list of question–answer segments. The schema lives in the `interview` package (issue #1).
- **Upload service**: a signed-in user uploads a report; the service validates it against the
  schema and writes its records to Postgres, with the author taken from the sign-in.
- Not in the MVP: the MCP server, the Mattermost bot, retrieval, transcription and private Sources.

## Rules

- Read `CONTEXT.md` for vocabulary, `docs/adr/` for decisions and `docs/pcr/` for conventions before
  changing a part they cover. Record a new decision as a new ADR or PCR, never by editing an accepted one.
- The backend is Go, with `go.mod` at the repo root; any SQL database is Postgres (PCR 0001, 0002).
- Skill output is untrusted. Validate every report in Go against its JSON Schema before anything is written
  to Postgres (PCR 0003). Schemas are `.json` files loaded with `go:embed`, never restated in Go (PCR 0004).
- Stored records are never updated or deleted. New information is a new record.
- Never return the author from an API, MCP tool or prompt that agents or the Board can reach.
- Interviews are personal data about named people. No real transcripts, reports or names in the repo,
  tests, fixtures or logs: use synthetic interviews and `@example.com`.
- PR titles and commits are Conventional Commits (`feat:`, `fix:`, `docs:`, …): they drive releases.

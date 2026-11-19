# AGENTS.md

agneta: captures what KTHAIS Advisors know but never wrote down, so the Board keeps that guidance after they
leave. The MVP only collects knowledge: an interview skill submits a JSON report through this repo's MCP
server, which validates it and stores its records in Postgres. Nothing reads them back yet.

## Commands

- Test: `make test` (vet and race-detected tests, as the PR check runs them)
- Run: `cp .env.example .env`, then `make run`. Every setting comes from the environment (PCR 0005).
- Image: `docker build -t agneta .`
- Database: `docker compose up -d` starts local Postgres; `go tool goose up` applies the migrations in
  `internal/db/migrations/` (PCR 0006). `docker compose down` stops it and keeps
  the data; `docker compose down -v` also deletes the data, for a fresh start.

## MVP

- **Interview skill**: an agent skill that runs an interview, gets the Advisor's approval and submits one
  report with the `submit_report` MCP tool.
- **Report**: one JSON file per interview: the interviewee's name, role and the date once at the top, then a
  list of question–answer segments. The schema lives in the `interview` package (issue #1).
- **MCP server**: this repo. `/mcp` serves `submit_report`, which validates the report against the schema
  (ADR 0001). Sign-in and Postgres storage are not built yet: a valid report is acknowledged with
  `stored: false`, and there is no author.
- Not in the MVP: the Mattermost bot, retrieval, transcription and private Sources.

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

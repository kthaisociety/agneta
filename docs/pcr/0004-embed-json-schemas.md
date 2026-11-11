# PCR-0004: Load JSON Schemas with go:embed

**Status:** Accepted

**Date:** 2026-11-09

## Convention

Every JSON Schema MUST live in the repo as a `.json` file and MUST be loaded into Go with `//go:embed`.

## Why

We weighed two other options.

- **Reading the file from disk at runtime.** The binary then fails if the container image leaves the file out, and the failure only shows at startup, not at build time.
- **Restating the schema as Go structs or a Go string.** That makes two copies, and they drift.

With `go:embed` the `.json` file is the one copy. The compiler puts it in the binary, so the tests, the ingestion code and the MCP tool all read the same bytes.

## How to apply

- Keep each schema as a plain `.json` file, next to the Go package that embeds it.
- Never paste a schema into a Go string or derive one from Go structs.

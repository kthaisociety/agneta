# PCR-0003: Validate documents against JSON Schema in Go

**Status:** Accepted

**Date:** 2026-11-09

## Convention

Every document an agent hands the backend MUST be validated in Go against a JSON Schema, draft 2020-12, using `github.com/santhosh-tekuri/jsonschema/v6` with format assertion turned on.

## Why

We weighed three other options.

- **Validating in Postgres with `pg_jsonschema`.** Managed Postgres hosts often don't offer the extension. A failed `CHECK` returns `false` with no field path, so the agent can't be told what to fix.
- **Go structs with strict decoding and hand-written checks.** The agent and its MCP tool need a schema file they can read, and hand-written checks drift from it.
- **`github.com/xeipuuv/gojsonschema`.** It stops at draft 7 and has had no maintenance for years.

`santhosh-tekuri/jsonschema/v6` supports draft 2020-12, passes the official JSON Schema test suite, and reports errors with the path to the bad field, such as `/segments/1/answer`.

## How to apply

- Turn on format assertion with `compiler.AssertFormat()`. Draft 2020-12 treats `format` as a note by default, so without it `"date": "2026-13-45"` passes.
- Close every object with `additionalProperties: false` unless the schema means to allow extra fields.
- Reject a document before anything writes it to Postgres.

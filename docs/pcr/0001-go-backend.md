# PCR-0001: Write the backend in Go

**Status:** Accepted

**Date:** 2026-11-09

## Convention

The agneta backend MUST be written in Go, and its tests MUST run with `go test`.

## Why

Go is the organisation's default backend language. A project moves off it only when Go lacks real support for something the project needs.

agneta doesn't hit that limit. Go builds one static binary, which fits the small container image the Dockerfile already uses, and the official MCP Go SDK covers the MCP tool the interview agent will call.

## How to apply

- Backend code lives in one Go module, with `go.mod` at the repo root.
- Tests sit next to the code they cover, in `_test.go` files.
- A backend file in another language needs a new PCR first, and that PCR has to name what Go can't support.

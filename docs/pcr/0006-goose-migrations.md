# PCR-0006: Change the schema only through goose migrations

**Status:** Proposed

**Date:** 2026-11-19

## Convention

Every schema change MUST be a goose SQL migration in `internal/db/migrations/`. Never edit a merged one.

## Why

goose keeps up and down in one SQL file, embeds with `embed.FS`, and is pinned in `go.mod`. Unlike
golang-migrate, **goose can also run a migration written in Go**, for a change SQL alone can't make, such as
re-embedding stored records with a new model. golang-migrate also splits each migration into two files.
Atlas adds a separate tool and a declarative workflow we don't need yet.

## How to apply

- Run it with `go tool goose`. Test `up`, `down`, `up` before a PR.

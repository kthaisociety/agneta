# PCR-0002: Use Postgres for any SQL database

**Status:** Accepted

**Date:** 2026-11-09

## Convention

Any SQL database agneta uses MUST be PostgreSQL.

## Why

Postgres is the organisation-wide SQL database, and it's the one our hosting supports. Another SQL engine would need hosting the organisation doesn't run.

## How to apply

- Interview data goes into Postgres tables.
- A different SQL engine, such as MySQL or SQLite, needs a new PCR first.
- Table constraints such as `NOT NULL` and `CHECK` are a backstop. Validating incoming documents is the backend's job, per [PCR-0003](0003-json-schema-validation-in-go.md).

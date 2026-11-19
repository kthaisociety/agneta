# 0002. Store reports in Postgres

Status: Proposed

## Context

ADR 0001 left storage open: `submit_report` validates a report, then drops it.

## Decision

- Tables `reports` and `records` (one per question–answer pair, ordered by `position`), named after
  `CONTEXT.md`. Roles in a `roles` table.
- Constraints are only a backstop; Go validates first (PCR-0003).
- Append-only: no cascading deletes.

## Consequences

- A new role changes both `schema.json` and `roles`.
- No author, summary or context columns yet. Nothing stops a resubmitted report being stored twice.

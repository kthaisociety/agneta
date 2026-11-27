# 0002. Store reports in Postgres

Status: Accepted

## Context

ADR 0001 left storage open: `submit_report` validates a report, then drops it.

## Decision

- Tables `reports` and `records` (one per question–answer pair, ordered by `position`), named after
  `CONTEXT.md`. `reports.role` is plain text: `schema.json` is the only list of roles.
- Constraints are only a backstop; Go validates first (PCR-0003).
- Append-only: no cascading deletes.

## Consequences

- A new role changes only `schema.json`. The database checks only that a role isn't blank.
- No author, summary or context columns yet. Nothing stops a resubmitted report being stored twice.

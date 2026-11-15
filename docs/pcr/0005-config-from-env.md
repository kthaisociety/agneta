# PCR-0005: All configuration comes from the environment, at runtime

**Status:** Accepted

**Date:** 2026-11-15

## Convention

Every setting MUST be an environment variable, read once in `internal/config` at startup. An invalid setting
stops the service at startup with one line naming it.

## Why

This is lucid's convention (its PCR 0001), and agneta deploys the same way: one image runs in staging and
production with different environments, set by `kthaisociety/deployments`. Build-time settings would differ
between the two. Secrets live in OpenBao, never in the repo or the image.

## How to apply

- A new setting means a field in `config.Config`, a line in `.env.example` and, when deploying, a line in
  the project's `project.yaml` in `kthaisociety/deployments`.
- `.env` is for local work only and is git-ignored. `cmd/agneta` loads it if present.
- Only the environment, never files or flags.

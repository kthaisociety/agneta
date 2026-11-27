# 0003. Sign in through identity-service, with agneta as Claude's authorization server

Status: Proposed

## Context
Anyone who can reach `/mcp` can call it (ADR 0001). Before agneta runs on a public host, it must:
- let in only people with a kthais.com account,
- know who submitted each report (the author), without ever showing it to agents or the Board,
- catch a missing or expired sign-in before an interview starts, not after 45 minutes of talking.

Two apps need a token, and no single sign-in gives both:
- **Claude's servers call `/mcp`**, so Claude must be an OAuth client of something that issues tokens for agneta.
- **agneta must know who the user is**, and only identity-service (IS) can say that. IS is the only
  KTHAIS service with Google credentials (IS ADR 0002), so agneta never talks to Google.

## Decision
agneta sits between two OAuth flows:

```
Claude ──(OAuth #1: agneta is the authorization server)──▶ agneta
agneta ──(OAuth #2: agneta is an IS client "agneta")─────▶ identity-service ──▶ Google
```

- **OAuth #1, Claude to agneta.** agneta issues its own tokens. It accepts exactly one client: Claude's
  published identity, client ID `https://claude.ai/oauth/mcp-oauth-client-metadata`, redirect URI
  `https://claude.ai/api/mcp/auth_callback`, public client (no secret), PKCE S256. Both values are
  constants. There is no dynamic client registration.
  - `/mcp` without a valid token answers `401` with `WWW-Authenticate` pointing to
    `/.well-known/oauth-protected-resource`.
  - The authorization server metadata advertises `client_id_metadata_document_supported: true` and the
    auth method `none`, so Claude uses its published identity.
  - Access tokens live 1 hour and refresh tokens rotate on every use. Both are random strings, stored only
    as hashes.
- **OAuth #2, agneta to IS.** It's the flow the IS team offered: client `agneta` with a secret, callback
  `PUBLIC_URL/auth/callback`, code + PKCE, and the token exchange over IS's internal URL.
  - agneta verifies the `id_token` (signature against IS's JWKS, `iss`, `aud=agneta`, `exp`) and reads
    `sub` and `email`.
- **Only `kthais.com`.** An email in any other domain gets a plain page asking the user to reconnect with
  their kthais.com account. Claude gets no code. The domain is a constant, not a setting.
- **A grant** links an agneta token to the person behind it: their IS user ID, email and IS refresh token.
  The author of a report comes from the grant, never from the report JSON.
- **Revocation happens at the hourly refresh.** When Claude refreshes its agneta token, agneta refreshes the
  IS token too.
  - If IS answers `invalid_grant`, agneta deletes the grant and answers `invalid_grant`, and Claude asks
    the Advisor to reconnect.
  - Any other IS failure keeps the grant and answers `temporarily_unavailable`.
- **`check_access`**, a new MCP tool, returns `{"authorized": true}`. The interview skill calls it first,
  so a broken connection shows up before the first question.
- **Grants and tokens live in memory** until the app connects to Postgres (ADR 0002 adds the database, but
  nothing in the app uses it yet).
- **`AUTH_DISABLED=true`** turns all of this off, for local work only.

## Consequences
- A restart or redeploy drops every grant, so every Advisor has to reconnect. That's acceptable for staging
  tests, and the grants move to Postgres once the app uses it for reports (ADR 0002, PCR-0002). The IS
  refresh token must be encrypted once it is stored there.
- Revoking someone at IS takes effect within an hour, at their next refresh. There is no cache and no
  `/oauth/revoke`.
- "Logging out" means removing the connector in Claude. An unused grant dies when its IS refresh token
  expires (30 days). agneta never holds an IS browser session, so IS's `/end_session` doesn't apply.
- Only claude.ai can connect. Claude Code or another MCP client would need its client ID added, which means
  a new decision.
- Running only one replica is required while grants live in memory.
- `AUTH_DISABLED` must never be set in staging or production. Nothing enforces this except review of the
  deployment config.
- Still open:
  - whether custom connectors work in Claude voice mode (to be tested)
  - agneta's public staging host (IT-88)
  - IS confirming that its token response includes an `id_token` with `email`

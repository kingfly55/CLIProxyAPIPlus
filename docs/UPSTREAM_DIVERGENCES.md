# Upstream Divergences

This working tree is intended to diverge from the upstream CLIProxyAPIPlus project in the following ways:

## 1. Plus-fork scope
- This fork is positioned as the "Plus" variant and is intended to carry community-maintained third-party provider support in addition to upstream functionality.
- The project README already documents that third-party provider support is the main reason this fork exists.

## 2. Custom management panel source
- The server configuration exposes `remote-management.panel-github-repository` so deployments can fetch a management panel from a forked panel repository instead of only the upstream panel source.
- For this environment, the intended usage is to serve a customized management center fork that includes enhanced usage analytics.

## 3. Local deployment customization
- The local deployment currently uses a customized `management.html` built from the forked `Cli-Proxy-API-Management-Center` repository.
- That customized panel adds auth-file-aware usage aggregation, per-model breakdowns, and client-side cost totals in the Usage page.

## Relevant files in this tree
- `README.md`
- `config.example.yaml`
- `internal/managementasset/updater.go`
- `internal/api/server.go`
- `internal/config/config.go`

## 4. Embedded models.json kept in sync with remote catalog
- `gpt-5.4-mini` was present in the remote models catalog (`router-for-me/models`) for all codex tiers but absent from the embedded `models.json` baked into the binary.
- Without this fix the model was only available after the background 3-hour refresh succeeded; on a fresh server start it would silently disappear until the first refresh.
- The embedded catalog is patched to include `gpt-5.4-mini` in `codex-free`, `codex-team`, `codex-plus`, and `codex-pro` tiers so the model is available immediately on startup.

## Relevant files in this tree
- `README.md`
- `config.example.yaml`
- `internal/managementasset/updater.go`
- `internal/api/server.go`
- `internal/config/config.go`
- `internal/registry/models/models.json`

## 5. Auth-file quota reset and capped quota retry recovery
- The local management API includes a manual quota reset action so auth files stuck in `quota_exceeded` can be made eligible again without waiting for the next automatic recovery path.
- The auth conductor also caps provider-supplied `Retry-After` values to the configured quota backoff ceiling before storing cooldown state, so abnormally large upstream retry values do not suspend a credential longer than the local backoff policy intends.
- This behavior exists specifically to support manual recovery of Codex/GPT-backed accounts in the local deployment and to keep quota cooling behavior predictable.

## Relevant files in this tree
- `internal/api/handlers/management/auth_files.go`
- `internal/api/server.go`
- `sdk/cliproxy/auth/conductor.go`

## 6. Suffixed OAuth model aliases can fork from base registered models
- The local deployment relies on `oauth-model-alias` entries whose source model name may include a thinking suffix, for example `gpt-5.4(xhigh)` backing the client-visible alias `claude-opus-4-6`.
- Upstream alias exposure logic only matched exact registered model IDs, which meant aliases backed by suffixed source names were loaded in config but not exposed in `/v1/models` and did not route at request time.
- The local fork extends alias exposure so a suffixed source name also matches the base registered model ID while preserving the forced suffix semantics for runtime routing. This keeps `claude-opus-4-6` visible and ensures it always maps to `gpt-5.4` with forced `xhigh` reasoning in this environment.

## Relevant files in this tree
- `sdk/cliproxy/service.go`
- `sdk/cliproxy/service_oauth_model_alias_test.go`
- `sdk/cliproxy/auth/oauth_model_alias.go`

## 7. Role-based cost attribution via original model slug tracking (KIN-40)
- WebCity role agents launch with `ANTHROPIC_DEFAULT_*_MODEL` env vars set to repo+role+tier slugs like `wa-builder-sonnet`, `mayor-opus`, `th-self-check-haiku`. These slugs are resolved to GPT-5.4 backing models via `oauth-model-alias.codex` entries in `config.yaml`.
- The fork captures the **original requested model slug** (as sent by the client, before any alias resolution) and threads it through the usage pipeline so every request can be attributed to its originating role, repo, and tier.
- New field `OriginalModel string` on `sdk/cliproxy/usage.Record` carries the slug from the handler to the usage plugin.
- New field `OriginalModel string` on `internal/usage.RequestDetail` (JSON `original_model`, `omitempty`) records the slug per individual request detail, so the frontend can filter/join by slug.
- `sdk/api/handlers/claude/code_handlers.go` sets `c.Set("original_model_slug", modelName)` in both streaming and non-streaming paths right after extracting the model from the request body (before any mapping happens).
- `internal/runtime/executor/helps/usage_helpers.go` introduces `OriginalModelFromContext(ctx)` (mirroring `APIKeyFromContext`) and `UsageReporter` threads the value through `buildRecord()`.
- `internal/usage/logger_plugin.go` adds an `OriginalSlugs map[string]*slugStats` aggregation bucket on each `apiStats`, and an `OriginalSlugs map[string]SlugSnapshot` field on `APISnapshot` (serialized as `original_slugs`, `omitempty`). `MergeSnapshot` merges these additively.
- New management endpoint `GET /v0/management/usage-by-role` (`internal/api/handlers/management/usage_by_role.go`) parses every slug into `(repo, role, tier)` via a known-prefix heuristic (`wa/cm/th`) and known tiers (`haiku/sonnet/opus`) and returns a per-role rollup with `by_repo`, `by_tier`, `by_slug` breakdowns plus an `unparsed` bucket. Wired in `internal/api/server.go` alongside the existing `/usage` routes.
- `disable-auto-update-panel: true` is set in deployment config (not code) to prevent the upstream auto-updater from overwriting the locally-built custom management panel. No new config option was needed — the existing flag suffices.

## Relevant files in this tree
- `sdk/cliproxy/usage/manager.go`
- `sdk/api/handlers/claude/code_handlers.go`
- `internal/runtime/executor/helps/usage_helpers.go`
- `internal/usage/logger_plugin.go`
- `internal/api/handlers/management/usage_by_role.go`
- `internal/api/handlers/management/usage_by_role_test.go`
- `internal/api/server.go`

## Notes
- The panel-specific functional divergence is committed in the separate forked management-center repository.
- This document records the intended differences between this local fork and the upstream `router-for-me/CLIProxyAPIPlus` repository.

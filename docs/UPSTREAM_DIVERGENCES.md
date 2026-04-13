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

## Notes
- The panel-specific functional divergence is committed in the separate forked management-center repository.
- This document records the intended differences between this local fork and the upstream `router-for-me/CLIProxyAPIPlus` repository.

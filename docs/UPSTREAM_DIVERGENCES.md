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

## Notes
- The panel-specific functional divergence is committed in the separate forked management-center repository.
- This document records the intended differences between this local fork and the upstream `router-for-me/CLIProxyAPIPlus` repository.

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

## Notes
- The panel-specific functional divergence is committed in the separate forked management-center repository.
- This document records the intended differences between this local fork and the upstream `router-for-me/CLIProxyAPIPlus` repository.

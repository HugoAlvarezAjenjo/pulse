# Roadmap

## v0.1.0 — Foundation ✅

- [x] Checks: `command`, `file`, `port`
- [x] YAML config with auto-discovery
- [x] Concurrent execution with deterministic output
- [x] Renderers: pretty (Lip Gloss) and plain
- [x] Fix system with confirmation prompts
- [x] Commands: `pulse`, `pulse check`, `pulse init`
- [x] Cross-platform support (Windows/macOS/Linux)
- [x] Unit tests
- [x] Documentation

## v0.2.0 — Distribution & CI ✅

- [x] GitHub Actions CI (Linux, macOS, Windows)
- [x] Goreleaser (automated binary releases)
- [x] GitHub Actions Release workflow
- [x] Install script (`curl | sh`)
- [x] Version info (`pulse --version`)
- [x] README badges

## v0.3.0 — Output Formats ✅

- [x] `--output json` structured output
- [x] `--output github` GitHub Actions annotations
- [x] `--quiet` only show failures

## v0.4.0 — New Check Types ✅

- [x] `env` — verify environment variables exist
- [x] `http` — health check endpoints
- [x] `docker` — verify containers are running

## v0.5.0 — Developer Experience ✅

- [x] `pulse doctor` — self-diagnostics
- [x] `pulse init` with auto-detection and `--preset`
- [x] `pulse init --empty` for blank templates
- [x] `pulse list` — inspect checks without running
- [x] `pulse ci` — CI-optimized mode with `--fail-fast`
- [x] Auto-detect TTY (plain output when piped)

## v0.6.0 — Teams & Scale ✅

- [x] Check groups (`--group` / `-g`)
- [x] `.pulse.local.yaml` personal overrides (merged automatically)
- [x] Per-check `timeout` in config
- [x] Global `--timeout` flag
- [x] Total execution duration in summary

## v1.0.0 — Stable Release

- [ ] Config `extends` (inheritance from base configs)
- [ ] Config schema validation with clear error messages
- [ ] `pulse validate` — check config syntax without running
- [ ] Update README with all current features
- [ ] Update architecture docs
- [ ] Man pages / CLI help improvements
- [ ] Semver stability guarantee

## v1.1.0 — Power Features

- [ ] `pulse watch` — re-run on config change
- [ ] Check `depends_on` (dependency between checks)
- [ ] `pulse report --markdown` — exportable reports
- [ ] Pre-commit hook integration guide

## v2.0.0 — Ecosystem

- [ ] Lightweight plugin system (external check binaries)
- [ ] Community preset registry
- [ ] VSCode extension
- [ ] Config registry for organizations

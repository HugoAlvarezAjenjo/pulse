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

## v0.3.0 — Output Formats

- [ ] `--output json` structured output
- [ ] `--output github` GitHub Actions annotations
- [ ] `--quiet` only show failures

## v0.4.0 — New Check Types

- [ ] `env` — verify environment variables exist
- [ ] `http` — health check endpoints
- [ ] `docker` — verify containers are running

## v0.5.0 — Developer Experience

- [ ] `pulse generate` — auto-detect stack and generate config
- [ ] `pulse doctor` — self-diagnostics
- [ ] Presets: `pulse init --preset node|go|python`

## v1.0.0 — Stable Release

- [ ] Config schema validation
- [ ] Check groups/profiles (`pulse check --group backend`)
- [ ] Global `--timeout` flag
- [ ] `pulse version` with build info
- [ ] Man pages
- [ ] Semver stability guarantee

## v1.1.0 — Power Features

- [ ] `pulse watch` — re-run on config change
- [ ] Config `extends` (inheritance)
- [ ] Check `depends_on`

## v1.2.0 — Team & Collaboration

- [ ] `.pulse.local.yaml` overrides
- [ ] Exportable markdown reports
- [ ] Pre-commit hook integration

## v2.0.0 — Ecosystem

- [ ] Lightweight plugin system
- [ ] Community preset registry
- [ ] VSCode extension
- [ ] `pulse ci` optimized mode

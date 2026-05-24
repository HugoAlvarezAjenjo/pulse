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
- [x] `.pulse.local.yaml` personal overrides
- [x] Per-check `timeout` in config
- [x] Global `--timeout` flag
- [x] Total execution duration in summary

## v1.0.0 — Stable Release ✅

- [x] Config `extends` (inheritance with circular detection)
- [x] `pulse validate` — check config syntax without running
- [x] `pulse report` — generate shareable markdown reports
- [x] `pulse completion` — shell auto-complete (bash/zsh/fish/powershell)
- [x] Full README with all features documented
- [x] Architecture docs updated
- [x] Example configs for 6 stacks (Node, Go, Python, Java, Monorepo, Extends)
- [x] Extends test suite

---

## v1.1.0 — Power Features

- [ ] `pulse watch` — re-run checks on config file change
- [ ] Check `depends_on` — dependency between checks
- [ ] `pulse diff` — compare environments between teammates
- [ ] Pre-commit hook integration guide + example

## v1.2.0 — Observability

- [ ] `pulse report --html` — HTML report with styling
- [ ] Onboarding metrics tracking (time to first green)
- [ ] Webhook/Slack notifications on first success
- [ ] `pulse status` — one-line summary for shell prompts

## v2.0.0 — Ecosystem

- [ ] Lightweight plugin system (external check binaries)
- [ ] Community preset registry
- [ ] VSCode extension (uses `--output json`)
- [ ] Config registry for organizations (remote `extends`)
- [ ] Landing page / documentation site

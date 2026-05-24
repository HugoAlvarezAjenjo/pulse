# Pulse

[![CI](https://github.com/HugoAlvarezAjenjo/pulse/actions/workflows/ci.yml/badge.svg)](https://github.com/HugoAlvarezAjenjo/pulse/actions/workflows/ci.yml)
[![Release](https://github.com/HugoAlvarezAjenjo/pulse/actions/workflows/release.yml/badge.svg)](https://github.com/HugoAlvarezAjenjo/pulse/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/HugoAlvarezAjenjo/pulse)](https://goreportcard.com/report/github.com/HugoAlvarezAjenjo/pulse)

> Fast CLI tool to validate development environments. Define checks in YAML, get clear diagnostics.

Pulse validates whether a development environment is correctly configured. It checks for required tools, files, services, and environment variables — then reports what's missing and optionally suggests fixes.

```
Pulse environment check

✓ Go go version go1.26.2 darwin/arm64
✓ Go module found go.mod
✓ Docker running
✗ PostgreSQL
  localhost:5432 refused connection

  Suggested fix:
  docker compose up db -d

Summary: 3 passed • 1 failed • 142ms
```

## Why Pulse?

- **Validate, don't provision** — checks your environment, never installs software
- **Declarative** — define checks in YAML, not imperative scripts
- **Fast** — concurrent execution, typically <200ms for 10+ checks
- **Zero config servers** — lives in your repo, runs locally
- **CI-native** — JSON output, GitHub annotations, proper exit codes
- **Team-friendly** — groups, local overrides, per-check timeouts

## Installation

### Quick install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/HugoAlvarezAjenjo/pulse/main/install.sh | sh
```

### Go install

```bash
go install github.com/HugoAlvarezAjenjo/pulse/cmd/pulse@latest
```

### Download binary

Pre-built binaries for Linux, macOS, and Windows are available on the [Releases](https://github.com/HugoAlvarezAjenjo/pulse/releases) page.

### Build from source

```bash
git clone https://github.com/HugoAlvarezAjenjo/pulse.git
cd pulse
make build
```

## Quick Start

```bash
# Auto-detect your stack and generate config
pulse init

# Or use a preset
pulse init --preset node

# Run checks
pulse
```

## Commands

| Command | Description |
|---------|-------------|
| `pulse` | Run all checks (same as `pulse check`) |
| `pulse check` | Run environment checks |
| `pulse ci` | Run checks in CI mode (no colors, no prompts) |
| `pulse init` | Generate a `.pulse.yaml` (auto-detect or preset) |
| `pulse list` | List configured checks without running them |
| `pulse doctor` | Show diagnostic info about pulse and config |

## Configuration

Create a `.pulse.yaml` in your project root:

```yaml
checks:
  - name: Node.js
    type: command
    command: node --version
    expected: ">=18"

  - name: Environment file
    type: file
    path: .env
    fix:
      run: "cp .env.example .env"

  - name: PostgreSQL
    type: port
    host: localhost
    port: 5432
    timeout: 5s
    groups: [backend]
    fix:
      run: "docker compose up db -d"

  - name: API Key
    type: env
    variable: API_KEY
    expected: "sk-*"
    groups: [backend]

  - name: Frontend deps
    type: command
    command: pnpm --version
    groups: [frontend]
```

## Check Types

### `command`

Validates a command exists and optionally checks its version.

```yaml
- name: Go
  type: command
  command: go version
  expected: ">=1.21"
```

Version extraction is automatic — works with any output format (`go version go1.26.2 darwin/arm64`, `node v22.3.0`, `Python 3.12.1`, etc.)

### `file`

Validates a file exists.

```yaml
- name: Config
  type: file
  path: .env
```

### `port`

Validates a TCP port is accepting connections.

```yaml
- name: Database
  type: port
  host: localhost
  port: 5432
```

### `env`

Validates an environment variable exists (optionally matches a pattern).

```yaml
- name: API Key
  type: env
  variable: API_KEY
  expected: "sk-*"    # wildcards: "sk-*", "*-production"
```

### `http`

Validates an HTTP endpoint responds with expected status.

```yaml
- name: API Health
  type: http
  url: http://localhost:3000/health
  status: 200
```

### `docker`

Validates a Docker container is running.

```yaml
- name: Redis
  type: docker
  container: redis
```

## Check Groups

Assign checks to groups for selective execution:

```yaml
checks:
  - name: Node.js
    type: command
    command: node --version
    groups: [frontend]

  - name: PostgreSQL
    type: port
    port: 5432
    groups: [backend, test]

  - name: .env file        # no group = always runs
    type: file
    path: .env
```

```bash
pulse check -g frontend     # only frontend + global checks
pulse check -g backend      # only backend + global checks
pulse check -g backend -g test   # union of groups
```

## Timeouts

```yaml
checks:
  - name: Slow build
    type: command
    command: make build
    timeout: 2m             # per-check timeout

  - name: Quick port
    type: port
    port: 5432
    timeout: 3s
```

```bash
pulse --timeout 10s         # global override for all checks
```

Priority: per-check > global flag > default (30s)

## Fix Commands

Any check can include a suggested fix:

```yaml
- name: Database
  type: port
  port: 5432
  fix:
    run: "docker compose up db -d"
```

```bash
pulse --fix    # prompts before executing each fix
```

## Local Overrides

Create `.pulse.local.yaml` for personal overrides (add to `.gitignore`):

```yaml
# .pulse.local.yaml — not committed
checks:
  - name: PostgreSQL       # overrides by name
    type: port
    host: localhost
    port: 5433             # my postgres is on a different port

  - name: Debug mode       # new check, only for me
    type: env
    variable: DEBUG
```

Merge rules:
- Same name → local overrides base completely
- New name → appended
- No groups removed from base

## Output Formats

```bash
pulse                    # pretty (auto-detects terminal)
pulse --plain            # no colors (auto when piped)
pulse -o json            # structured JSON
pulse -o github          # GitHub Actions annotations
pulse -q                 # quiet: only show failures
```

Auto-TTY detection: colors are disabled automatically when output is piped or redirected.

## CI Integration

### Basic

```bash
pulse ci                    # plain, quiet, no prompts
pulse ci --fail-fast        # stop at first failure
pulse ci -o github          # GitHub Actions annotations
pulse ci -g backend         # only backend checks
```

### GitHub Actions

```yaml
- name: Validate environment
  run: pulse ci --output github
```

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | One or more checks failed |
| 2 | Config/runtime error |

## All Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `pretty`, `plain`, `json`, `github` |
| `--quiet` | `-q` | Only show failures and errors |
| `--plain` | | Disable colors (shorthand for `-o plain`) |
| `--group` | `-g` | Run only specified groups (repeatable) |
| `--timeout` | `-t` | Timeout per check (e.g., `10s`, `1m`) |
| `--config` | | Path to config file |
| `--fix` | | Prompt to run suggested fixes |
| `--fail-fast` | | Stop at first failure (ci command) |
| `--version` | | Show version info |
| `--preset` | | Use preset template (init command) |
| `--empty` | | Generate blank template (init command) |
| `--force` | | Overwrite existing config (init command) |

## Config Inheritance (`extends`)

Share a base config across multiple projects:

```yaml
# base.pulse.yaml (shared across org)
checks:
  - name: Git
    type: command
    command: git --version
  - name: Docker
    type: command
    command: docker --version
```

```yaml
# .pulse.yaml (project-specific)
extends: base.pulse.yaml

checks:
  - name: Node.js
    type: command
    command: node --version
    expected: ">=20"
```

Rules:
- Child checks with same `name` override the parent
- New checks are appended
- Supports chains (parent can also `extends`)
- Max 5 levels deep (circular detection included)
- Path is relative to the config file

## Config Discovery

Pulse searches in the current directory:

1. `.pulse.yaml` (preferred)
2. `pulse.yaml`
3. Resolves `extends` chain if present
4. `.pulse.local.yaml` (merged on top, if exists)

## Development

```bash
make build      # Build binary
make test       # Run all tests
make lint       # Run golangci-lint
make fmt        # Format code
make run        # Build and run
make install    # Install to $GOPATH/bin
```

## Architecture

See [docs/architecture.md](docs/architecture.md) for design details.

## License

MIT

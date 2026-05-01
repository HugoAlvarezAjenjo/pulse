# Pulse

> Takes the pulse of your development environment.

Pulse validates whether a developer environment is correctly configured for a project. It checks for required tools, files, and services — then reports what's missing and optionally suggests fixes.

## Philosophy

- **Validate, don't provision** — Pulse checks your environment, it doesn't install software
- **Declarative configuration** — Define checks in YAML, not scripts
- **Fast and minimal** — No background daemons, no TUI, no bloat
- **Clear diagnostics** — When something fails, you know why and how to fix it

## Installation

```bash
go install github.com/HugoAlvarezAjenjo/pulse/cmd/pulse@latest
```

Or build from source:

```bash
git clone https://github.com/HugoAlvarezAjenjo/pulse.git
cd pulse
make build
```

## Quick Start

1. Create a `.pulse.yaml` in your project root:

```yaml
checks:
  - name: Node.js
    type: command
    command: node --version
    expected: ">=22"

  - name: Environment file
    type: file
    path: .env

  - name: PostgreSQL
    type: port
    host: localhost
    port: 5432
    fix:
      run: "docker compose up db -d"
```

2. Run pulse:

```bash
pulse
```

Output:

```
Pulse environment check

✓ Node.js v22.3.0
✓ Environment file found .env
✗ PostgreSQL
  localhost:5432 refused connection
  ensure a service is listening on localhost:5432

  Suggested fix:
  docker compose up db -d

Summary: 2 passed • 1 failed
```

## Commands

### `pulse` / `pulse check`

Run all environment checks defined in the configuration.

```bash
pulse              # same as pulse check
pulse check        # explicit
pulse --fix        # prompt to run fixes for failed checks
pulse --plain      # disable colors (for CI)
pulse --config path/to/config.yaml
```

### `pulse init`

Generate a starter `.pulse.yaml` in the current directory.

```bash
pulse init
```

## Configuration Reference

### Check Types

#### `command`

Validates that a command exists and optionally checks its version output.

```yaml
- name: Node.js
  type: command
  command: node --version
  expected: ">=22"
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Display name for the check |
| `type` | yes | Must be `command` |
| `command` | yes | Command to execute |
| `expected` | no | Version constraint (supports `>=`) |
| `fix` | no | Fix command to suggest |

#### `file`

Validates that a file exists at the specified path.

```yaml
- name: Environment file
  type: file
  path: .env
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Display name |
| `type` | yes | Must be `file` |
| `path` | yes | Path to check (relative to cwd) |
| `fix` | no | Fix command to suggest |

#### `port`

Validates that a TCP port is accepting connections.

```yaml
- name: PostgreSQL
  type: port
  host: localhost
  port: 5432
  fix:
    run: "docker compose up db -d"
```

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Display name |
| `type` | yes | Must be `port` |
| `host` | no | Host to check (default: `localhost`) |
| `port` | yes | Port number |
| `fix` | no | Fix command to suggest |

### Fix Commands

Any check can include a `fix` block with a `run` field:

```yaml
fix:
  run: "docker compose up db -d"
```

When `--fix` is passed and a check fails, pulse will:
1. Display the suggested fix command
2. Ask for confirmation
3. Execute the command
4. Stream output in real time

## Flags

| Flag | Description |
|------|-------------|
| `--plain` | Disable colors and styling |
| `--config` | Path to config file (overrides auto-discovery) |
| `--fix` | Prompt to run suggested fixes |

## Config Discovery

Pulse looks for configuration in the current directory:

1. `.pulse.yaml` (preferred)
2. `pulse.yaml`

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | One or more checks failed |
| 2 | Internal/config/runtime error |

## Development

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linter
make fmt        # Format code
make install    # Install to $GOPATH/bin
```

## Architecture

See [docs/architecture.md](docs/architecture.md) for design details.

## License

MIT

# Architecture

## Overview

Pulse follows clean architecture principles with clear separation of concerns. The codebase is organized into focused packages with minimal coupling.

```
cmd/pulse/          → Entry point (main.go)
cli/                → CLI commands (Cobra)
  check.go          → Run environment checks
  ci.go             → CI-optimized mode
  init.go           → Generate config (auto-detect/presets)
  list.go           → List checks without running
  validate.go       → Validate config syntax
  report.go         → Generate markdown reports
  doctor.go         → Self-diagnostics
  completion.go     → Shell completions
  root.go           → Root command + global flags
internal/
  checks/           → Check interface + implementations
  config/           → Config parsing, discovery, extends, merge
  generator/        → Stack detection + preset templates
  renderer/         → Output rendering (pretty, plain, json, github)
  runner/           → Concurrent check execution
  fix/              → Fix command execution
  result/           → Result model
  styles/           → Lip Gloss style definitions
```

## Design Principles

### Interfaces Over Concrete Types

Core abstractions are defined as interfaces:

- **`checks.Check`** — any validation that can be run
- **`renderer.Renderer`** — any way to output results

### Composition Over Inheritance

Each check type is a simple struct implementing the `Check` interface. No base class, no embedding hierarchy.

### Deterministic Output

Checks run concurrently but results are always rendered in config order.

## Data Flow

```
Config Discovery
    ↓
Extends Resolution (parent → child chain)
    ↓
Local Merge (.pulse.local.yaml)
    ↓
Group Filter (--group flag)
    ↓
Factory (CheckConfig → Check implementations)
    ↓
Runner (concurrent execution with timeouts)
    ↓
Results (ordered []Result)
    ↓
Renderer (pretty/plain/json/github)
    ↓
Exit Code (0/1/2)
```

## Config Resolution Order

```
1. Discover .pulse.yaml or pulse.yaml
2. If `extends` → resolve parent chain (max 5 levels, circular detection)
3. Merge: child overrides parent by check name
4. If .pulse.local.yaml exists → merge on top
5. Final config passed to runner
```

## Check Types

| Type | What it validates |
|------|-------------------|
| `command` | Command exists + optional version constraint |
| `file` | File exists at path |
| `port` | TCP port accepting connections |
| `env` | Environment variable set + optional pattern |
| `http` | HTTP endpoint responds with expected status |
| `docker` | Docker container is running |

## Concurrency Model

The runner launches one goroutine per check. Results are stored in a pre-allocated slice indexed by position. A `sync.WaitGroup` coordinates completion.

Each check gets its own `context.Context` with timeout:
- Per-check timeout (config) > Global timeout (flag) > Default (30s)

## Renderer System

```go
type Renderer interface {
    Header()
    Success(r result.Result)
    Failure(r result.Result)
    Error(r result.Result)
    Summary(results []result.Result, duration time.Duration)
}
```

Implementations:
- **PrettyRenderer** — Colored terminal output (Lip Gloss), auto-disabled when piped
- **PlainRenderer** — No colors, for CI/pipes
- **JSONRenderer** — Structured JSON for tooling integration
- **GitHubRenderer** — GitHub Actions workflow commands

## Error Model

| Status | Meaning | Exit Code |
|--------|---------|-----------|
| Success | Check passed | 0 |
| Failure | Expected condition not met | 1 |
| Error | Runtime/internal error | 2 |

## Package Dependencies

```
cmd/pulse → cli → internal/*
                    ├── checks → config, result
                    ├── config → (yaml.v3)
                    ├── generator → (os, filepath)
                    ├── runner → checks, result
                    ├── renderer → result, styles
                    ├── fix → result
                    ├── result → (none)
                    └── styles → (lipgloss)
```

The `result` package has zero internal dependencies. The `config` package only depends on yaml.v3. Clean, acyclic dependency graph.

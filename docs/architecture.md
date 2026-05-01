# Architecture

## Overview

Pulse follows clean architecture principles with clear separation of concerns. The codebase is organized into focused packages with minimal coupling between them.

```
cmd/pulse/          → Entry point
cli/                → CLI commands (Cobra)
internal/
  checks/           → Check interface + implementations
  config/           → Configuration parsing + discovery
  renderer/         → Output rendering interface + implementations
  runner/           → Concurrent check execution
  fix/              → Fix command execution
  result/           → Result model
  styles/           → Lip Gloss style definitions
```

## Design Principles

### Interfaces Over Concrete Types

Core abstractions are defined as interfaces:

- **`checks.Check`** — Any validation that can be run
- **`renderer.Renderer`** — Any way to output results

This enables testing, extensibility, and future output formats (JSON, JUnit XML, etc.) without modifying existing code.

### Composition Over Inheritance

Each check type is a simple struct implementing the `Check` interface. There's no base class, no embedding hierarchy. Each check owns its logic completely.

### Deterministic Output

Although checks run concurrently for speed, results are always rendered in the order they were defined in configuration. This ensures reproducible, predictable output regardless of execution timing.

## Data Flow

```
Config File → Parse → []CheckConfig → Factory → []Check → Runner → []Result → Renderer → Output
```

1. **Discovery** — Find `.pulse.yaml` or `pulse.yaml`
2. **Parsing** — Decode YAML into typed config structs
3. **Factory** — Convert config entries into Check implementations
4. **Runner** — Execute all checks concurrently with timeouts
5. **Rendering** — Output results using the selected renderer
6. **Fix** — Optionally prompt and execute fix commands

## Concurrency Model

The runner launches one goroutine per check. Results are stored in a pre-allocated slice indexed by position, ensuring order preservation without sorting. A `sync.WaitGroup` coordinates completion.

Each check gets its own `context.Context` with a timeout, enabling:
- Individual check timeouts
- Global cancellation
- Clean shutdown

## Error Model

Three distinct result states:

| Status | Meaning | Exit Code |
|--------|---------|-----------|
| Success | Check passed | 0 |
| Failure | Check didn't pass (expected condition not met) | 1 |
| Error | Runtime/internal error (command crashed, invalid config) | 2 |

This distinction is important: a port being closed is a **failure** (the environment isn't ready), but a malformed config is an **error** (the tool can't function).

## Renderer System

The `Renderer` interface decouples output formatting from check execution:

```go
type Renderer interface {
    Header()
    Success(r result.Result)
    Failure(r result.Result)
    Error(r result.Result)
    Summary(results []result.Result)
}
```

Current implementations:
- **PrettyRenderer** — Colored terminal output with Lip Gloss
- **PlainRenderer** — Unformatted text for CI/pipes

Future implementations could include:
- JSON output
- JUnit XML (for CI integration)
- GitHub Actions annotations

## Fix System

The fix system is intentionally simple:
- Checks can declare a fix command in config
- Pulse never auto-executes fixes
- User must confirm each fix interactively
- Commands run with a timeout
- Output is streamed in real time

This keeps pulse as a **diagnostic tool**, not a provisioning system.

## Package Dependencies

```
cmd/pulse → cli → internal/*
                    ├── checks → config, result
                    ├── config → (yaml.v3)
                    ├── runner → checks, result
                    ├── renderer → result, styles
                    ├── fix → result
                    ├── result → (none)
                    └── styles → (lipgloss)
```

The `result` package has zero internal dependencies, making it a stable foundation. The `config` package only depends on the YAML library. This keeps the dependency graph clean and acyclic.

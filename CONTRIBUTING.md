# Contributing to Pulse

Thank you for your interest in contributing to Pulse!

## Development Setup

```bash
git clone https://github.com/HugoAlvarezAjenjo/pulse.git
cd pulse
go mod tidy
make build
```

## Running Tests

```bash
make test        # Run all tests
make lint        # Run linter
```

## Branch Strategy

We use **GitHub Flow**:

1. Create a branch from `main`: `feat/your-feature` or `fix/your-fix`
2. Make your changes with clear, focused commits
3. Open a Pull Request against `main`
4. CI must pass before merge
5. Squash merge when approved

## Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add env check type
fix: handle Windows paths in file check
docs: update configuration reference
test: add port check timeout tests
chore: update CI workflow
```

## Code Guidelines

- Follow idiomatic Go conventions
- Keep functions small and focused
- Add tests for new functionality
- Run `make lint` before submitting
- No dependencies without strong justification

## Adding a New Check Type

1. Create `internal/checks/yourcheck.go` implementing the `Check` interface
2. Add config fields to `internal/config/config.go`
3. Register in `internal/checks/factory.go`
4. Add tests in `internal/checks/yourcheck_test.go`
5. Update documentation

## Philosophy

Before contributing, please understand what Pulse is and isn't:

**Pulse does:**
- Validate environments
- Report problems clearly
- Execute user-defined fixes (with confirmation)

**Pulse does NOT:**
- Install software automatically
- Become a package manager
- Add cloud integrations
- Create TUI interfaces

Keep it minimal, fast, and trustworthy.

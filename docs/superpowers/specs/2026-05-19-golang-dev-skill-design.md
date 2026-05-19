# golang-dev Skill Design

## Overview

A Go development best-practices guide skill that serves as a **background reference** when writing Go code. It covers CLI scaffolding, configuration management, common libraries, build/test commands, and escape analysis verification.

**Type:** Background guide + user-invocable (`/golang-dev`)
**Trigger:** When writing Go code or starting a new Go project
**Mode:** Read-write (can generate/modify code)

## Skill Metadata

```yaml
name: golang-dev
description: >
    Go development best-practices guide covering CLI scaffolding (cobra),
    configuration (viper), common libraries (zap, testify, golangci-lint),
    build/test commands, and escape analysis verification. Use when writing
    Go code, starting a new Go project, or choosing libraries for Go development.
allowed-tools: Read, Grep, Glob, Bash, Edit, Write
user-invocable: true
disable-model-invocation: false
```

## Content Structure (7 Sections)

### 1. CLI Development with cobra

**When to use:** Any Go project that exposes a command-line interface.

**Content:**

- Always use `github.com/spf13/cobra` for CLI applications
- Project structure convention: `cmd/` directory with one file per subcommand
- Root command setup pattern: `rootCmd` in `cmd/root.go`
- Subcommand registration via `init()` or explicit `AddCommand()`
- Flag binding pattern: `cobra.Command.Flags()` for command-specific, `PersistentFlags()` for inherited
- Integration with viper for flag-to-config binding via `viper.BindPFlag()`

**Code example:** Root command + subcommand scaffold showing `RunE` (not `Run`) for error propagation.

### 2. Configuration with viper

**When to use:** Any Go project that needs configuration management.

**Content:**

- Always use `github.com/spf13/viper` for configuration loading
- Loading precedence (highest to lowest): environment variables > CLI flags > config file > defaults
- Supported config formats: YAML (preferred), TOML, JSON
- Config file search: `viper.AddConfigPath()` with multiple paths (`.`, `$HOME`, `/etc/app/`)
- Environment variable binding: `viper.AutomaticEnv()` + `viper.SetEnvPrefix()`
- cobra integration: `viper.BindPFlag()` in `init()` to sync flags and config
- Struct unmarshalling: `viper.Unmarshal(&cfg)` into a typed config struct

**Code example:** Full config loading setup with cobra integration.

### 3. Common Libraries

Categorized library recommendations for Go development:

| Category   | Library                       | When to use                                                                                                  |
| ---------- | ----------------------------- | ------------------------------------------------------------------------------------------------------------ |
| Logging    | `go.uber.org/zap`             | Production-grade structured logging. Use `zap.NewProduction()` for JSON, `zap.NewDevelopment()` for console. |
| Testing    | `github.com/stretchr/testify` | Assertions (`assert`), requirements (`require`), and mocks (`mock`) for tests.                               |
| Linting    | `golangci-lint`               | Meta-linter aggregating 50+ linters. Run via `golangci-lint run ./...`.                                      |
| HTTP       | `net/http` (stdlib)           | Prefer stdlib for most cases. Use `github.com/gin-gonic/gin` only when routing complexity justifies it.      |
| Hot Reload | `github.com/air-verse/air`    | Development-time hot reload. Run `air` instead of `go run`.                                                  |
| Mocking    | `github.com/bytedance/mockey` | Runtime monkey-patching for tests (ByteDance internal).                                                      |
| CLI        | `github.com/spf13/cobra`      | CLI scaffolding (see Section 1).                                                                             |
| Config     | `github.com/spf13/viper`      | Configuration management (see Section 2).                                                                    |

### 4. Build Commands

**Common build commands and flags:**

- Standard build: `go build ./...`
- Production binary (stripped): `go build -ldflags="-s -w" -o bin/app ./cmd/app`
    - `-s` strips symbol table
    - `-w` strips DWARF debug info
    - Result: ~30-40% smaller binary
- Version injection: `go build -ldflags="-X main.version=1.2.3 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"`
- Static binary (no CGO): `CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app ./cmd/app`
- Cross-compilation: `GOOS=linux GOARCH=amd64 go build -o bin/app-linux ./cmd/app`
- Race detector build: `go build -race ./...`

### 5. Test Commands

**Common test commands and flags:**

- Run all tests: `go test ./...`
- Verbose output: `go test -v ./...`
- Race detector: `go test -race ./...`
- Coverage report: `go test -cover -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out`
- Specific test: `go test -run TestFunctionName ./pkg/...`
- Benchmarks: `go test -bench=. -benchmem -run=^$ ./...`
- Short mode (skip slow tests): `go test -short ./...`
- Disable inlining and optimizations: `go test -gcflags="all=-N -l" ./...`
    - `-N` disables compiler optimizations
    - `-l` disables inlining
    - `all=` applies flags to all compiled packages, not just the test package
    - Use cases: debugging with `dlv`, accurate escape analysis, diagnosing optimization-related issues
- Count (disable test caching): `go test -count=1 ./...`
- Timeout: `go test -timeout 30s ./...`

### 6. Escape Analysis Verification

**When to use:** Investigating memory allocations, verifying stack vs heap allocation decisions.

**Commands:**

- Basic escape analysis: `go build -gcflags='-m' ./...`
- Detailed escape analysis: `go build -gcflags='-m=2' ./...`
- Specific package: `go build -gcflags='-m=2' ./pkg/handler/...`

**Common escape reasons and fixes:**

| Escape reason                         | Typical cause                             | Fix                                                |
| ------------------------------------- | ----------------------------------------- | -------------------------------------------------- |
| `leaking param`                       | Parameter stored beyond function scope    | Avoid storing pointer params in long-lived structs |
| `moved to heap: too large`            | Stack frame exceeds limit (~10MB)         | Break into smaller allocations or use `sync.Pool`  |
| `moved to heap: captured by closure`  | Variable captured in a closure            | Pass value explicitly instead of capturing         |
| `moved to heap: interface conversion` | Assigning concrete to `interface{}`/`any` | Use concrete types in hot paths                    |
| `&x escapes to heap`                  | Returning address of local variable       | Return value instead of pointer when possible      |

**Verification workflow:**

1. Run `go build -gcflags='-m=2' ./... 2>&1 | grep "escapes to heap"` to find escaping allocations
2. Focus on hot paths (handlers, loops, frequently called functions)
3. Use `go test -bench=. -benchmem` to measure allocation impact
4. Apply fixes, re-run escape analysis to verify

### 7. Quick Reference Table

A summary table of all commands for fast lookup:

| Task             | Command                                                        |
| ---------------- | -------------------------------------------------------------- |
| Build (dev)      | `go build ./...`                                               |
| Build (prod)     | `go build -ldflags="-s -w" -o bin/app ./cmd/app`               |
| Build (static)   | `CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app ./cmd/app` |
| Test             | `go test ./...`                                                |
| Test (race)      | `go test -race ./...`                                          |
| Test (cover)     | `go test -cover -coverprofile=coverage.out ./...`              |
| Test (no inline) | `go test -gcflags="all=-N -l" ./...`                           |
| Test (bench)     | `go test -bench=. -benchmem -run=^$ ./...`                     |
| Escape analysis  | `go build -gcflags='-m=2' ./...`                               |
| Lint             | `golangci-lint run ./...`                                      |
| Vet              | `go vet ./...`                                                 |
| Format           | `gofmt -w .` or `goimports -w .`                               |
| Hot reload       | `air`                                                          |
| Cross compile    | `GOOS=linux GOARCH=amd64 go build -o bin/app ./cmd/app`        |

## Relationship with Existing Skills

- **golang-code-quality**: Reviews code quality (SOLID, architecture). `golang-dev` complements it by guiding library choices and build workflow.
- **golang-performance-tuning**: Deep performance review. `golang-dev` covers escape analysis basics; performance-tuning goes deeper.
- **golang-mvc**: Architecture patterns. `golang-dev` doesn't overlap; it covers tooling and workflow.
- **golang-naming**: Naming conventions. No overlap.
- **golang-dead-code**: Dead code removal. No overlap.
- **golang-network**: Network performance. No overlap.

## File Location

`skills/golang-dev/SKILL.md` in the cc-plugin project.

# ADR-001: Why Go for the Runtime

> Status: Accepted
> Date: 2026-07-10

## Context

NAEOS needs a primary implementation language for the runtime engine, pipeline orchestrator, and CLI. Key requirements:

- Fast compilation for rapid development cycles
- Excellent concurrency support for parallel pipeline execution and multi-target code generation
- Small, statically-linked binaries for easy distribution
- Strong standard library for networking, file I/O, and JSON/YAML processing
- Cross-compilation to all major platforms

## Decision

Use **Go** as the primary language for the NAEOS runtime.

## Consequences

### Positive

- Single binary distribution with no runtime dependencies
- Goroutines provide lightweight concurrency for pipeline steps and streaming
- Go modules provide reproducible builds
- Mature ecosystem for CLI tools (`cobra`, `viper`) and API servers (`net/http`, `chi`)
- Native WASM compilation via `GOOS=wasip1 GOARCH=wasm` for browser and plugin sandboxing

### Negative

- WASM target is needed for non-Go plugins that must run in sandboxed environments
- Go generics are limited compared to other languages, requiring more verbose code in some areas
- Runtime reflection has overhead; compile-time code generation is preferred for hot paths

### Mitigations

- Plugin SDK uses WASM for third-party extensions, keeping the core runtime in native Go
- Code generation via `go generate` and NEIR consumes NEIR to produce type-safe Go code
- Benchmark tooling (see `naeos benchmark`) tracks performance regressions

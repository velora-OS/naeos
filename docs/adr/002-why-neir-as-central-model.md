# ADR-002: Why NEIR as the Central Model

> Status: Accepted
> Date: 2026-07-10

## Context

NAEOS generates code for multiple target languages (Go, TypeScript, Python, Java, Rust) from a single specification. Each generator needs a consistent, language-agnostic intermediate representation to consume. Without a shared IR, every spec-to-language path requires a dedicated compiler, leading to combinatorial complexity.

## Decision

Use **NEIR** (NAEOS Engine Intermediate Representation) as the central intermediate model that all generators consume.

## Consequences

### Positive

- `N` generators require `N` NEIR-to-language compilers instead of `N×M` spec-to-language compilers
- Validation, optimization, and transformation happen once at the NEIR level
- Source maps from spec language constructs to NEIR nodes enable precise error reporting
- NEIR serializes to JSON/YAML for inspection and debugging (`naeos inspect`)

### Negative

- All generators must understand NEIR, adding a dependency on NEIR schema stability
- NEIR schema changes require coordinated updates across all generators
- Debugging generation failures requires understanding the NEIR intermediate state

### Mitigations

- NEIR schema is versioned (see NES-023) and follows semver
- `naeos validate` checks NEIR integrity before generation
- Generator tests consume NEIR fixtures to catch schema drift early
- The `naeos benchmark` command tracks compilation throughput across NEIR versions

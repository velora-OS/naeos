# NES-033 — Testing Guide

> Status: Draft
> Last Updated: 2026-07-10

Guide for running and writing tests in the NAEOS codebase.

---

## Running Tests

### All Tests

```bash
go test ./...
```

### Specific Package

```bash
go test ./pkg/kernel/...
go test ./internal/specification/parser/...
go test ./cmd/naeos/...
```

### Verbose Output

```bash
go test -v ./...
```

### With Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Race Detector

```bash
go test -race ./...
```

---

## Test Structure

Tests use the standard `testing` package. No external test frameworks.

### Basic Pattern

```go
func TestSomething(t *testing.T) {
    result := DoSomething()
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Table-Driven Tests

```go
func TestCompare(t *testing.T) {
    tests := []struct {
        a, b    string
        want    int
    }{
        {"a", "b", -1},
        {"b", "a", 1},
        {"a", "a", 0},
    }
    for _, tt := range tests {
        if got := Compare(tt.a, tt.b); got != tt.want {
            t.Errorf("Compare(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
        }
    }
}
```

### Temporary Directories

```go
func TestWritesFile(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "output.txt")
    os.WriteFile(path, []byte("content"), 0644)
    // ... verify file
}
```

### Error Checking

```go
func TestFailsOnInvalidInput(t *testing.T) {
    _, err := Parse("{invalid")
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}
```

---

## Test Files by Package

### pkg/kernel

| File | Tests |
|---|---|
| `kernel_test.go` | Register/resolve, start/stop, event bus, telemetry |

### pkg/config

| File | Tests |
|---|---|
| `config_test.go` | JSON loading, YAML loading |

### pkg/pipeline

| File | Tests |
|---|---|
| `pipeline_test.go` | Run produces result, injected parser, output dir, config loading |

### internal/specification/parser

| File | Tests |
|---|---|
| `parser_test.go` | JSON parse, YAML parse, invalid input, structured spec, full spec |

### internal/specification/normalizer

| File | Tests |
|---|---|
| `normalizer_test.go` | Convert parsed spec, architecture/deployment/testing, nil error |

### internal/specification/resolver

| File | Tests |
|---|---|
| `resolver_test.go` | Build context from normalized spec |

### internal/neir/builder

| File | Tests |
|---|---|
| `builder_test.go` | Build NEIR, extract architecture, nil input |

### internal/neir/validator

| File | Tests |
|---|---|
| `validator_test.go` | Valid/incomplete/nil NEIR, duplicates, port validation |

### internal/neir/version

| File | Tests |
|---|---|
| `version_test.go` | SemVer parsing, comparison, compatibility, validation |

### internal/planner/graph

| File | Tests |
|---|---|
| `graph_test.go` | Nodes, edges, topological sort, cycle detection |

### internal/planner/scheduler

| File | Tests |
|---|---|
| `scheduler_test.go` | Task generation, no-services, fallback |

### internal/registry

| File | Tests |
|---|---|
| `registry_test.go` | Register, resolve, unregister, find by category/version |

### internal/events

| File | Tests |
|---|---|
| `bus_test.go` | Pub/sub, multiple subscribers, unsubscribe, topics |

### internal/generation/engine

| File | Tests |
|---|---|
| `engine_test.go` | Generate artifacts from NEIR |

### internal/generation/renderers

| File | Tests |
|---|---|
| `renderer_test.go` | Templates, custom functions, nested data |

### internal/governance/policy

| File | Tests |
|---|---|
| `evaluator_test.go` | All condition operators, disabled rules, default rules |

### internal/governance/review

| File | Tests |
|---|---|
| `reviewer_test.go` | Review rules: no-todo, no-placeholder, package decl |

### internal/runtime/engine

| File | Tests |
|---|---|
| `engine_test.go` | Execute, validate, history, reset, batch execution |

### internal/runtime/telemetry

| File | Tests |
|---|---|
| `telemetry_test.go` | Emit, metrics, filter by name, reset |

### internal/knowledge/graph

| File | Tests |
|---|---|
| `graph_test.go` | Nodes, edges, find by type/metadata, path checking |

### internal/knowledge/provenance

| File | Tests |
|---|---|
| `provenance_test.go` | Record, lineage, find by artifact/source/creator |

### cmd/naeos

| File | Tests |
|---|---|
| `main_test.go` | All CLI commands, output formats, flags |

---

## Writing New Tests

### 1. Create Test File

Place `*_test.go` in the same package:

```go
// internal/myPackage/myfile_test.go
package myPackage

import "testing"

func TestMyFunction(t *testing.T) {
    // ...
}
```

### 2. Use Subtests

```go
func TestParse(t *testing.T) {
    t.Run("valid input", func(t *testing.T) {
        // ...
    })
    t.Run("empty input", func(t *testing.T) {
        // ...
    })
}
```

### 3. Use Helper Functions

```go
func setupTest(t *testing.T) *MyStruct {
    t.Helper()
    return NewMyStruct()
}
```

### 4. Use Temp Directories

```go
func TestWriteFile(t *testing.T) {
    dir := t.TempDir()
    // Files auto-cleanup after test
}
```

---

## Test Conventions

| Convention | Description |
|---|---|
| File naming | `*_test.go` in same package |
| Function naming | `TestFunctionName` |
| Subtest naming | `t.Run("description", ...)` |
| Temp dirs | `t.TempDir()` for auto-cleanup |
| Error checking | `if err == nil { t.Fatal(...) }` |
| No external deps | Use standard `testing` package only |
| Table-driven | Use `[]struct{...}` pattern for multiple cases |

---

## CI Integration

Tests run automatically via GitHub Actions:

```yaml
# .github/workflows/ci.yml
- name: Run tests
  run: go test ./...
```

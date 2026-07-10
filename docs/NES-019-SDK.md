# NES-019 SDK

## 1. Status
- Status: Draft
- Version: 0.3
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the programmatic integration layer for building tools and plugins on top of NAEOS, including multi-language SDK generation.

## 3. Scope
The SDK covers public API packages (Go), multi-language output adapters, integration patterns, and extension development.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The SDK shall expose public API packages under `pkg/`.
- FR-002: The SDK shall provide constructors and interfaces for all public components.
- FR-003: The SDK shall support custom pipeline configurations.
- FR-004: The SDK shall generate projects in 5 languages: Go, TypeScript, Python, Java, Rust.
- FR-005: The SDK shall use the OutputAdapter pattern for language extensibility.

### 4.2 Non-Functional Requirements
- NFR-001: Public APIs shall be stable and documented.
- NFR-002: SDK usage shall not require knowledge of internal packages.
- NFR-003: Each generated project shall build and run without modification.

## 5. Public Packages

| Package | Deskripsi |
|---------|-----------|
| `pkg/kernel` | Core kernel: service registry, lifecycle, event bus, telemetry |
| `pkg/pipeline` | Pipeline orchestrator |
| `pkg/config` | Configuration loader |

### 5.1 pkg/kernel

```go
k := kernel.NewKernel()

// Service registry
k.Register("my-service", myService)

// Lifecycle
k.Initialize()
k.Start()
defer k.Stop()

// Event bus
k.Publish("topic", payload)
k.Subscribe("topic", handler)

// Telemetry
k.RecordEvent(kernel.TelemetryEvent{Name: "event", Timestamp: time.Now().Unix()})
metrics := k.Metrics()
```

### 5.2 pkg/pipeline

```go
p := pipeline.NewPipeline(config)

// Full run
results, err := p.Run("specification.yaml")

// Validate only
results, err = p.Validate("specification.yaml")
```

### 5.3 pkg/config

```go
cfg, err := config.Load("config.yaml")
// cfg.Name, cfg.Mode, cfg.Verbose, cfg.OutputDir
```

## 6. Extension Points

| Extension | Deskripsi |
|-----------|-----------|
| Custom Generator | Tambahkan target bahasa baru |
| Custom Validator | Tambahkan aturan validasi baru |
| Custom Policy | Tambahkan aturan governance baru |
| Custom Renderer | Tambahkan engine rendering baru |

## 7. Integration Patterns

### 7.1 Embed NAEOS in Your Application

```go
import "github.com/NAEOS-foundation/naeos/pkg/kernel"

k := kernel.NewKernel()
// Register your services
// Use kernel lifecycle
```

### 7.2 Build Custom Pipeline

```go
import "github.com/NAEOS-foundation/naeos/pkg/pipeline"

p := pipeline.NewPipeline(customConfig)
results, err := p.Run("my-spec.yaml")
```

## 8. Multi-Language SDK Generation

NAEOS dapat menghasilkan proyek dalam berbagai bahasa dari satu spesifikasi.

### 8.1 Supported Languages

| Language | Build File | Package Manager | Dockerfile |
|----------|------------|-----------------|------------|
| Go | go.mod | go mod | golang:1.22-alpine |
| TypeScript | package.json | npm | node:22-alpine |
| Python | pyproject.toml | pip | python:3.12-slim |
| Java | pom.xml | maven | eclipse-temurin:21 |
| Rust | Cargo.toml | cargo | rust:1.78-alpine |

### 8.2 Configuration

```yaml
pipeline:
  name: my-project
  language:
    - go
    - typescript
    - python
```

### 8.3 OutputAdapter Interface

```go
type OutputAdapter interface {
    Language() language.Language
    GenerateProject(projectName string) []Artifact
    GenerateModule(moduleName, modulePath, projectName string) []Artifact
    GenerateService(serviceName, serviceKind string, servicePort int, projectName string) []Artifact
    GenerateDockerfile(projectName string) []Artifact
    GenerateCI(projectName string) []Artifact
    GenerateDockerCompose(projectName string) []Artifact
    GenerateArchitectureDoc(projectName, pattern string) []Artifact
}
```

### 8.4 Adding New Languages

Implement `OutputAdapter` interface and register via `init()`:

```go
type MyAdapter struct{}

func init() { Register(MyAdapter{}) }
func (MyAdapter) Language() language.Language { return "mylang" }
// ... implement all methods
```

See [NES-039 SDK Multi-Language](NES-039-SDK-MultiLanguage.md) and [NES-040 Output Adapter Architecture](NES-040-Output-Adapter-Architecture.md) for full details.

## 9. Acceptance Criteria
- Public APIs are stable and documented.
- SDK usage does not require internal package knowledge.
- Extension points support custom generators, validators, and renderers.
- Integration patterns work for embedding and standalone usage.
- Generated projects build and run in all 5 supported languages.
- New languages can be added by implementing the OutputAdapter interface.

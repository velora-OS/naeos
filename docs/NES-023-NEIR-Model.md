# NAEOS NEIR Model Reference

## Status
- Status: Stable
- Version: 1.0.0
- Owner: NAEOS Foundation
- Last Updated: 2026-07-10

---

## 1. Overview

NEIR (Nusantara Engineering Intermediate Representation) adalah model engineering sentral yang merepresentasikan seluruh aspek sistem secara independen dari syntax atau format transport. Document ini mendeskripsikan struktur model aktual dari implementasi Go.

---

## 2. Package Structure

```
internal/neir/
├── model/
│   ├── model.go           # Main NEIR struct
│   ├── ai/                # AI model
│   ├── api/               # API model
│   ├── architecture/      # Architecture model
│   ├── component/         # Component model
│   ├── deployment/        # Deployment model
│   ├── domain/            # Domain model
│   ├── docs/              # Documentation model
│   ├── infrastructure/    # Infrastructure model
│   ├── metadata/          # Metadata model
│   ├── module/            # Module model
│   ├── project/           # Project model
│   ├── security/          # Security model
│   ├── service/           # Service model
│   ├── storage/           # Storage model
│   └── testing/           # Testing model
├── builder/               # Build NEIR from resolved spec
├── validator/             # Validate NEIR model
└── version/               # Version parsing and comparison
```

---

## 3. Core Model

### 3.1 NEIR Struct

```go
type NEIR struct {
    Project        *project.Project
    Architecture   *architecture.Architecture
    Domain         *domain.Domain
    Modules        []module.Module
    Components     []component.Component
    Services       []service.Service
    APIs           []api.API
    Storage        []storage.Storage
    Infrastructure *infrastructure.Infrastructure
    Security       *security.Security
    AI             *ai.AI
    Documentation  *docs.Documentation
    Deployment     *deployment.Deployment
    Testing        *testingmodel.Testing
    Metadata       *metadata.Metadata
}
```

---

## 4. Domain Models

### 4.1 Project

```go
type Project struct {
    Name        string
    Version     string
    Description string
    License     string
    Authors     []string
    Repository  string
    Tags        []string
    Attributes  map[string]string
}
```

**Fields:**
- `Name`: Project name (required)
- `Version`: Semantic version string
- `Description`: Project description
- `License`: License identifier (e.g., "MIT", "Apache-2.0")
- `Authors`: List of author names
- `Repository`: Git repository URL
- `Tags`: Classification tags
- `Attributes`: Custom key-value pairs

---

### 4.2 Architecture

```go
type Pattern string

const (
    PatternLayered    Pattern = "layered"
    PatternClean      Pattern = "clean"
    PatternHexagonal  Pattern = "hexagonal"
    PatternMicrokernel Pattern = "microkernel"
    PatternEventDriven Pattern = "event-driven"
    PatternCQRS       Pattern = "cqrs"
    PatternMonolith   Pattern = "monolith"
)

type Architecture struct {
    Pattern     Pattern
    Style       string
    Description string
    Principles  []string
    Layers      []Layer
    Attributes  map[string]string
}

type Layer struct {
    Name        string
    Description string
    Modules     []string
}
```

**Supported Patterns:**
- `layered`: Traditional N-tier architecture
- `clean`: Clean Architecture (Uncle Bob)
- `hexagonal`: Ports and Adapters
- `microkernel`: Plugin-based architecture
- `event-driven`: Event-driven architecture
- `cqrs`: Command Query Responsibility Segregation
- `monolith`: Single deployable unit

---

### 4.3 Module

```go
type Module struct {
    Name         string
    Path         string
    Description  string
    Packages     []string
    Dependencies []string
    Attributes   map[string]string
}
```

**Fields:**
- `Name`: Module name (required)
- `Path`: File system path
- `Description`: Module description
- `Packages`: Go packages in this module
- `Dependencies`: Other module dependencies
- `Attributes`: Custom key-value pairs

---

### 4.4 Service

```go
type ServiceKind string

const (
    KindHTTP   ServiceKind = "http"
    KindGRPC   ServiceKind = "grpc"
    KindWorker ServiceKind = "worker"
    KindCLI    ServiceKind = "cli"
    KindJob    ServiceKind = "job"
)

type Service struct {
    Name        string
    Kind        ServiceKind
    Port        int
    Description string
    Endpoints   []Endpoint
    Middleware  []string
    Attributes  map[string]string
}

type Endpoint struct {
    Method string
    Path   string
    Action string
}
```

**Service Kinds:**
- `http`: HTTP/REST service
- `grpc`: gRPC service
- `worker`: Background worker
- `cli`: Command-line tool
- `job`: Batch job

---

### 4.5 Storage

```go
type StorageType string

const (
    TypeSQL   StorageType = "sql"
    TypeNoSQL StorageType = "nosql"
    TypeFile  StorageType = "file"
    TypeCache StorageType = "cache"
    TypeQueue StorageType = "queue"
    TypeBlob  StorageType = "blob"
)

type Storage struct {
    Name        string
    Type        StorageType
    Provider    string
    Connection  string
    Collections []Collection
    Attributes  map[string]string
}

type Collection struct {
    Name   string
    Schema map[string]string
}
```

---

### 4.6 Security

```go
type Security struct {
    Authentication *Authentication
    Authorization  *Authorization
    Encryption     *Encryption
    Secrets        []Secret
    Attributes     map[string]string
}

type Authentication struct {
    Method   string
    Provider string
}

type Authorization struct {
    Model string
    Roles []string
}

type Encryption struct {
    InTransit bool
    AtRest    bool
    Algorithm string
}

type Secret struct {
    Name string
    Kind string
}
```

---

### 4.7 Deployment

```go
type Strategy string

const (
    StrategyRolling   Strategy = "rolling"
    StrategyBlueGreen Strategy = "blue-green"
    StrategyCanary    Strategy = "canary"
    StrategyRecreate  Strategy = "recreate"
)

type Deployment struct {
    Target       string
    Strategy     Strategy
    Environments []Environment
    Scaling      *Scaling
    Attributes   map[string]string
}

type Environment struct {
    Name   string
    Kind   string
    Config map[string]string
}

type Scaling struct {
    Min      int
    Max      int
    Replicas int
}
```

---

### 4.8 Testing

```go
type TestingStrategy string

const (
    StrategyUnit        TestingStrategy = "unit"
    StrategyIntegration TestingStrategy = "integration"
    StrategyE2E         TestingStrategy = "e2e"
    StrategyContract    TestingStrategy = "contract"
)

type Testing struct {
    Strategy   TestingStrategy
    Frameworks []string
    Coverage   *Coverage
    Fixtures   []Fixture
    Attributes map[string]string
}

type Coverage struct {
    MinPercent float64
}

type Fixture struct {
    Name string
    Kind string
    Path string
}
```

---

### 4.9 Infrastructure

```go
type Provider string

const (
    ProviderAWS   Provider = "aws"
    ProviderGCP   Provider = "gcp"
    ProviderAzure Provider = "azure"
    ProviderLocal Provider = "local"
)

type Infrastructure struct {
    Provider   Provider
    Region     string
    Resources  []Resource
    Networking []Network
    Attributes map[string]string
}

type Resource struct {
    Name string
    Kind string
    Spec map[string]string
}

type Network struct {
    Name  string
    Kind  string
    Ports []int
}
```

---

## 5. Builder

### 5.1 Builder Interface

```go
type Builder interface {
    Build(resolved any) (*model.NEIR, error)
}
```

The builder transforms resolved specification into a complete NEIR model.

### 5.2 Default Builder

Extracts:
- Project name and description
- Architecture pattern and principles
- Modules with paths and dependencies
- Services with kinds and ports
- Dependencies between modules

---

## 6. Validator

### 6.1 Validation Rules

| Rule | Severity | Description |
|------|----------|-------------|
| project-required | error | Project name must be present |
| project-name-valid | error | Project name must be non-empty |
| modules-required | error | At least one module required |
| module-name-valid | error | Module name must be non-empty |
| module-path-valid | error | Module path must be non-empty |
| services-valid | warning | Services should have valid configuration |
| service-port-valid | warning | Port should be in valid range (1-65535) |
| service-kind-valid | warning | Service kind should be recognized |
| architecture-pattern-valid | warning | Architecture pattern should be recognized |
| duplicate-modules | warning | Duplicate module names detected |
| duplicate-services | warning | Duplicate service names detected |

### 6.2 ValidateDetailed

```go
func (v DefaultValidator) ValidateDetailed(neir *model.NEIR) (*ValidationResult, error)
```

Returns detailed validation result with multiple errors.

```go
type ValidationResult struct {
    Valid  bool
    Errors []ValidationError
}

type ValidationError struct {
    Field   string
    Message string
    Severity string
}
```

---

## 7. Version

### 7.1 SemVer

```go
type Version struct {
    Major int
    Minor int
    Patch int
    Pre   string
}

func Parse(version string) (*Version, error)
func (v *Version) Compare(other *Version) int
func (v *Version) IsCompatible(constraint string) bool
func Validate(version string) error
```

---

## 8. Example NEIR

```json
{
  "project": {
    "name": "acme-api",
    "version": "1.0.0",
    "description": "ACME API service",
    "license": "MIT"
  },
  "architecture": {
    "pattern": "hexagonal",
    "description": "Clean architecture with ports and adapters",
    "principles": ["separation of concerns", "dependency inversion"]
  },
  "modules": [
    {
      "name": "auth",
      "path": "./internal/auth",
      "description": "Authentication module",
      "dependencies": ["crypto", "storage"]
    },
    {
      "name": "users",
      "path": "./internal/users",
      "description": "User management module"
    }
  ],
  "services": [
    {
      "name": "gateway",
      "kind": "http",
      "port": 8080,
      "description": "API gateway",
      "endpoints": [
        {"method": "GET", "path": "/api/v1/health", "action": "healthCheck"},
        {"method": "POST", "path": "/api/v1/auth/login", "action": "login"}
      ]
    }
  ],
  "deployment": {
    "strategy": "rolling",
    "environments": [
      {"name": "development"},
      {"name": "staging"},
      {"name": "production"}
    ]
  },
  "testing": {
    "strategy": "unit-integration",
    "coverage": {"min_percent": 80.0}
  }
}
```

---

## 9. Pipeline Integration

```
Specification (YAML/JSON)
    │
    ▼
┌─────────────┐
│   Parser    │  → SpecDocument
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Normalizer │  → NormalizedSpec
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Resolver   │  → ResolvedSpec
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Builder   │  → NEIR
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Validator  │  → ValidationResult
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Generator  │  → Artifacts
└─────────────┘
```

---

## 10. References

- [NES-023-NEIR.md](NES-023-NEIR.md) - NEIR Specification
- [NES-026-Pipeline.md](NES-026-Pipeline.md) - Pipeline Documentation
- [NES-002-Kernel-API.md](NES-002-Kernel-API.md) - Kernel API Reference

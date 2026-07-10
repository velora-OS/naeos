# NES-006 Template

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the template model — reusable patterns for generating consistent artifacts across NAEOS projects.

## 3. Scope
The template model covers template definition, variable substitution, conditional rendering, and template composition.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: Templates shall support variable substitution.
- FR-002: Templates shall support conditional sections.
- FR-003: Templates shall be composable (template inheritance).
- FR-004: Templates shall produce consistent output for equivalent inputs.

### 4.2 Non-Functional Requirements
- NFR-001: Templates shall be human-readable and maintainable.
- NFR-002: Template rendering shall be deterministic.

## 5. Template Model

### 5.1 Rendering Engine

Templates are rendered using Go's `text/template` engine via the `internal/generation/renderers` package.

```go
// Basic rendering
result, err := renderer.Render(templateString, data)

// Named rendering
result, err := renderer.RenderNamed("config", templateString, data)

// With custom functions
result, err := renderers.RenderWithFuncs("config", templateString, data, funcs)
```

### 5.2 Template Data

Standard data structure used across all templates:

```go
type TemplateData struct {
    Project    string            // Project name
    Module     string            // Module name
    Service    string            // Service name
    Port       int               // Service port
    Kind       string            // Artifact kind (http, grpc)
    Version    string            // Artifact version
    Package    string            // Go package name
    Attributes map[string]string // Custom attributes
}
```

### 5.3 Template Categories

#### Project Templates
- `main.go` — application entry point
- `go.mod` — Go module definition
- `Dockerfile` — container build
- `README.md` — project documentation
- CI/CD workflows

#### Module Templates
- Domain model (`domain.go`)
- Repository interface and implementation (`repository.go`)
- Service layer (`service.go`)
- HTTP handlers (`handler.go`)
- Router (`router.go`)
- Middleware (`middleware.go`)
- Configuration (`config.go`)
- Tests (`*_test.go`)

#### Service Templates
- Service configuration
- Endpoint definitions
- Health check handlers

## 6. Workflow
1. Select appropriate template based on artifact type.
2. Populate template data from NEIR model.
3. Apply conditional sections based on configuration.
4. Render template to produce output artifact.
5. Validate rendered output.

## 7. Built-in Functions

| Function | Description |
|----------|-------------|
| `toUpper` | Convert string to uppercase |
| `toLower` | Convert string to lowercase |
| `indent` | Add indentation |
| `join` | Join string slice with separator |

## 8. Acceptance Criteria
- Templates produce consistent output for equivalent inputs.
- Variable substitution works correctly for all data types.
- Conditional sections render based on template data.
- Custom functions are available in template rendering.

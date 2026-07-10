# NES-003 Workspace

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the workspace as the execution context for a NAEOS project.

## 3. Scope
This document covers workspace state, project configuration, local artifacts, and interaction with CLI or SDK tools.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The workspace shall maintain a consistent project state.
- FR-002: The workspace shall support local generation and validation of artifacts.
- FR-003: The workspace shall provide a root directory for all project artifacts.
- FR-004: The workspace shall support loading configuration from YAML or JSON files.
- FR-005: The workspace shall isolate generated artifacts from source specifications.

### 4.2 Non-Functional Requirements
- NFR-001: Workspace execution shall be reproducible across supported environments.
- NFR-002: Workspace state shall be inspectable and auditable.
- NFR-003: Workspace operations shall be idempotent where possible.

## 5. Workspace Model

### 5.1 Structure

```
project-root/
├── specification.yaml       # NAEOS specification file
├── config.yaml              # Pipeline configuration (optional)
├── config.json              # Pipeline configuration (optional)
├── out/                     # Generated artifacts
│   ├── <project-name>/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   └── <module>/
│   │   ├── go.mod
│   │   ├── Dockerfile
│   │   └── README.md
│   └── ...
└── .naeos/                  # Internal state (future)
```

### 5.2 Components

#### Project Metadata
Derived from the specification file. Includes project name, modules, services, and architecture settings.

#### Configuration Descriptors
Pipeline configuration loaded from `config.yaml` or `config.json`. Controls execution mode, output directory, and verbosity.

#### Generated Artifacts
Output files produced by the generator engine. Placed in the output directory specified by configuration.

#### Execution Cache
Internal state for tracking which pipeline stages have been executed (planned for future versions).

#### Dependency Manifests
Generated `go.mod` files for Go projects, or equivalent dependency files for other target languages.

## 6. Workflow

1. **Initialize** the workspace by locating the specification file.
2. **Load** the active configuration (YAML or JSON).
3. **Parse** the specification into internal representation.
4. **Execute** the requested operation (run, validate, inspect, etc.).
5. **Persist** generated artifacts and diagnostics to the output directory.

## 7. Configuration

### 7.1 Pipeline Configuration

```yaml
pipeline:
  name: my-project
  mode: full          # full | validate-only
  verbose: true
  output_dir: out
```

### 7.2 Supported Formats

- YAML (`.yaml`, `.yml`)
- JSON (`.json`)

Configuration loader tries JSON first, falls back to YAML if not found.

## 8. Acceptance Criteria
- A project can be initialized into a consistent workspace state.
- Workspace operations can be repeated without manual state repair.
- Configuration can be loaded from both YAML and JSON formats.
- Generated artifacts are placed in the correct output directory.

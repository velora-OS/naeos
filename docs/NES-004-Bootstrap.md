# NES-004 Bootstrap

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team
- Last Updated: 2026-07-09

## 2. Purpose
This specification defines the bootstrap procedure used to initialize a NAEOS project from initial input to a ready workspace. Bootstrap encompasses both the CLI `init`/`scaffold` commands for new project creation and the pipeline initialization sequence that prepares the runtime for specification processing.

## 3. Scope
The bootstrap process covers:
- Project initialization via `naeos init` (config file generation)
- Project scaffolding via `naeos scaffold` (full workspace creation)
- Pipeline construction and kernel service registration
- Initial validation of workspace structure
- Default configuration generation

### 3.1 Out of Scope
- Specification parsing and normalization (NES-030)
- Code generation and artifact output (NES-007)
- Governance policy evaluation (NES-027)

## 4. Definitions
| Term | Definition |
|------|-----------|
| **Bootstrap** | The process of initializing a NAEOS project workspace or pipeline from zero state |
| **Scaffold** | A complete project workspace generated from a project name and template |
| **Config File** | A YAML or JSON file containing pipeline configuration (name, mode, output_dir, language) |
| **Pipeline Config** | The in-memory `pipeline.Config` struct loaded from a config file or CLI flags |
| **Workspace** | The directory tree containing all project files, modules, and configuration |

## 5. Bootstrap Modes

NAEOS provides two distinct bootstrap modes:

### 5.1 Config Bootstrap (`naeos init`)
Generates a minimal pipeline configuration file.

```
User → naeos init -o config.yaml → config.yaml
```

**Output**: A YAML file with default pipeline settings:
```yaml
pipeline:
  name: naeos-dev
  mode: development
  verbose: true
  output_dir: ./out
```

### 5.2 Scaffold Bootstrap (`naeos scaffold`)
Generates a complete project workspace with all standard files.

```
User → naeos scaffold --name myapp --output ./myapp → full directory tree
```

**Output directory structure**:
```
myapp/
├── README.md
├── spec.yaml
├── Makefile
├── .gitignore
├── Dockerfile
├── .github/workflows/ci.yml
├── go.mod
├── cmd/app/main.go
├── config.yaml
├── config.json
└── internal/core/
    ├── README.md
    ├── package.go
    ├── config.yaml
    ├── handler.go
    ├── repository.go
    ├── service.go
    ├── domain/model.go
    ├── http/handler.go
    ├── http/router.go
    ├── middleware/logging.go
    ├── config/config.go
    ├── config/load.go
    └── handler_test.go
```

### 5.3 Pipeline Bootstrap (Runtime)
Constructs the in-memory pipeline from config and prepares it for execution.

```
Config File → pipeline.ConfigFromFile() → pipeline.New(cfg) → Kernel.Start() → Ready
```

## 6. Inputs

| Bootstrap Mode | Inputs |
|---------------|--------|
| Config init | Output file path (default: `config.example.yaml`) |
| Scaffold | Project name (required), output directory (default: project name) |
| Pipeline | Config file path (YAML/JSON), optional `--language` overrides |

## 7. Outputs

| Bootstrap Mode | Outputs |
|---------------|---------|
| Config init | YAML config file at specified path |
| Scaffold | Complete workspace directory tree with 20+ files |
| Pipeline | Initialized `*Pipeline` with all kernel services registered |

## 8. Requirements

### 8.1 Functional Requirements
- **FR-001**: The system shall create a canonical workspace structure for a new project via `naeos scaffold`.
- **FR-002**: The system shall initialize project metadata and initial configuration including `spec.yaml`, `go.mod`, `Dockerfile`, and CI workflow.
- **FR-003**: The system shall generate a minimal pipeline config via `naeos init`.
- **FR-004**: The pipeline shall register all core services (parser, normalizer, resolver, builder, validator, scheduler, generator, graph, registry, evaluator, reviewer) with the kernel during construction.
- **FR-005**: The pipeline shall start the kernel lifecycle (calling Initialize/Start on Lifecycle services) before executing any pipeline operation.
- **FR-006**: The pipeline shall stop the kernel lifecycle (calling Stop on Lifecycle services) after completing any pipeline operation.
- **FR-007**: The scaffold shall generate a working Go module with proper `go.mod`, entry point, and core module files.
- **FR-008**: The scaffold shall include a starter `spec.yaml` that can be immediately processed by `naeos run`.
- **FR-009**: The scaffold shall generate a Dockerfile, CI workflow, and standard project files (README, .gitignore, Makefile).

### 8.2 Non-Functional Requirements
- **NFR-001**: Bootstrap shall be repeatable and deterministic — identical inputs produce identical outputs.
- **NFR-002**: Bootstrap shall fail safely and report validation issues clearly (missing project name, invalid output path).
- **NFR-003**: Config file generation shall produce valid YAML that can be parsed by the config loader.
- **NFR-004**: Scaffold shall not overwrite existing files unless the output directory is clean.

## 9. Pipeline Construction Sequence

The pipeline construction follows this exact sequence:

```
1. ConfigFromFile(path)
   └→ LoadFile(path) → parse JSON/YAML → File{Pipeline{Name, Mode, Verbose, OutputDir, Language}}

2. pipeline.New(cfg)
   ├→ Create Pipeline struct with all field assignments
   ├→ Default-initialize nil components:
   │   ├→ parser.NewParser()
   │   ├→ normalizer.NewNormalizer()
   │   ├→ resolver.NewResolver()
   │   ├→ builder.NewBuilder()
   │   ├→ validator.NewValidator()
   │   ├→ scheduler.NewScheduler()
   │   ├→ engine.NewEngine()
   │   ├→ graph.New()
   │   ├→ registry.NewRegistry()
   │   ├→ policy.NewEvaluator()
   │   ├→ review.NewReviewer()
   │   └→ kernel.NewKernel()
   └→ registerKernelServices()
       └→ kernel.Register() for each service (parser, normalizer, resolver, etc.)

3. Pipeline.Run(input) / Pipeline.Validate(input)
   ├→ kernel.Start()
   │   └→ lifecycle.Initialize() + lifecycle.Start() for each Lifecycle service
   ├→ executeWithKernel(fn)
   │   ├→ emitKernelEvent("kernel.start", ...)
   │   ├→ fn() — the actual pipeline work
   │   └→ kernel.Stop()
   │       └→ lifecycle.Stop() for each Lifecycle service
   └→ Return result
```

## 10. Workflow

### 10.1 Config Init Workflow
1. Parse CLI flags (`--output` for file path).
2. Generate YAML content with default pipeline settings.
3. Write file to specified path.
4. Print confirmation: `created <path>`.

### 10.2 Scaffold Workflow
1. Validate required `--name` flag is provided.
2. Determine output directory (defaults to project name).
3. Create output directory structure via `os.MkdirAll`.
4. Generate and write all standard files (README, spec, Makefile, Dockerfile, CI, go.mod, main.go).
5. Generate and write all core module files (handler, service, repository, domain model, HTTP handler, middleware, config).
6. Print confirmation: `scaffolded <output>`.

### 10.3 Pipeline Bootstrap Workflow
1. Load config from file via `ConfigFromFile(path)`.
2. Merge CLI `--language` overrides into `cfg.Languages`.
3. Construct pipeline via `pipeline.New(cfg)`.
4. Kernel registers all 11+ services during construction.
5. On first Run/Validate call, kernel.Start() initializes all Lifecycle services.
6. Pipeline execution occurs within `executeWithKernel` wrapper.
7. After execution, kernel.Stop() shuts down all Lifecycle services.

## 11. Scaffold File Details

| File | Purpose |
|------|---------|
| `README.md` | Project overview with quick start instructions |
| `spec.yaml` | Starter NAEOS specification with project name, core module, and API service |
| `Makefile` | Build targets for help and scaffold |
| `.gitignore` | Excludes /bin/, /out/, *.log |
| `Dockerfile` | Go 1.22 Alpine build image |
| `.github/workflows/ci.yml` | GitHub Actions CI with checkout, setup-go, test |
| `go.mod` | Go module declaration |
| `cmd/app/main.go` | Application entrypoint with HTTP server |
| `config.yaml` / `config.json` | Starter configuration files |
| `internal/core/*.go` | Core module: handler, service, repository, domain model, HTTP layer, middleware, config |

## 12. Acceptance Criteria
- **AC-001**: A new project can be initialized using a single `naeos scaffold --name <name>` command.
- **AC-002**: The resulting workspace is valid for subsequent `naeos run` or `naeos validate` operations.
- **AC-003**: `naeos init` produces a valid YAML config file parseable by `pipeline.ConfigFromFile()`.
- **AC-004**: Pipeline construction registers all required services and they are available via kernel resolution.
- **AC-005**: Kernel lifecycle (Start/Stop) is properly invoked around every pipeline execution.
- **AC-006**: Bootstrap is deterministic — repeated scaffolds with the same name produce identical output.

## 13. Implementation References
- `cmd/naeos/main.go:72-98` — `newInitCommand()` implementation
- `cmd/naeos/main.go:389-466` — `newScaffoldCommand()` implementation
- `pkg/pipeline/pipeline.go:72-83` — `ConfigFromFile()` config loading
- `pkg/pipeline/pipeline.go:85-143` — `New()` pipeline construction with defaults
- `pkg/pipeline/pipeline.go:145-167` — `registerKernelServices()` service registration
- `pkg/pipeline/pipeline.go:169-188` — `executeWithKernel()` lifecycle wrapper
- `pkg/config/config.go` — Config file parsing (JSON/YAML)

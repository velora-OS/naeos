# NES-040 Output Adapter Architecture

## 1. Status
- Status: Draft
- Version: 0.1
- Owner: NAEOS Core Team

## 2. Purpose
Dokumentasi arsitektur Output Adapter — mekanisme extensible untuk menghasilkan artefak dalam berbagai bahasa pemrograman.

## 3. Scope
Dokumen ini mencakup interface adapter, registry, dispatch mechanism, dan pola extensi.

## 4. Normative References
- NES-039 SDK Multi-Language
- NES-013 Compiler
- NES-007 Generator

## 5. Architecture

### 5.1 Adapter Registry Pattern

```
OutputAdapter Interface
    ↓
┌──────────────────────────────────────────┐
│           Adapter Registry               │
│  map[language.Language]OutputAdapter     │
├──────────┬───────────┬──────────┬────────┤
│GoAdapter │TSAdapter  │PyAdapter │...     │
│init()    │init()     │init()    │        │
│Register()│Register() │Register()│        │
└──────────┴───────────┴──────────┴────────┘
```

Setiap adapter di-register secara otomatis melalui `init()` function saat package di-import.

### 5.2 Dispatch Flow

```
1. NEIR Model diterima
2. GenerationConfig.Languages diekstrak
3. Untuk setiap bahasa:
   a. Adapter diambil dari registry
   b. Method-method dipanggil secara sequential
   c. Artefak dikumpulkan
4. Semua artefak dikembalikan
```

## 6. Interface

```go
type OutputAdapter interface {
    Language() language.Language
    GenerateProject(projectName string) []engine.Artifact
    GenerateModule(moduleName, modulePath, projectName string) []engine.Artifact
    GenerateService(serviceName, serviceKind string, servicePort int, projectName string) []engine.Artifact
    GenerateDockerfile(projectName string) []engine.Artifact
    GenerateCI(projectName string) []engine.Artifact
    GenerateDockerCompose(projectName string) []engine.Artifact
    GenerateArchitectureDoc(projectName, pattern string) []engine.Artifact
}
```

## 7. Registry Operations

```go
// Register an adapter (called in init())
func Register(adapter OutputAdapter)

// Get adapter for a language
func Get(lang language.Language) (OutputAdapter, bool)

// Get all registered adapters
func All() map[language.Language]OutputAdapter

// Generate artifacts for a NEIR model using configured languages
func GenerateForNEIR(neir *model.NEIR) ([]engine.Artifact, error)
```

## 8. Language Package

Package `internal/neir/model/language` menyediakan metadata per bahasa:

```go
func IsValid(lang Language) bool           // Cek apakah bahasa didukung
func All() []Language                     // Semua bahasa yang didukung
func Extensions(lang Language) []string    // File extensions per bahasa
func BuildFile(lang Language) string       // Build file per bahasa
func DockerBaseImage(lang Language) string // Docker base image
func DockerRuntimeImage(lang Language) string // Docker runtime image
```

## 9. Adapter Implementation Pattern

```go
type MyAdapter struct{}

func init() {
    Register(MyAdapter{})
}

func (MyAdapter) Language() language.Language {
    return language.LanguageMyLang
}

func (MyAdapter) GenerateProject(projectName string) []engine.Artifact {
    return []engine.Artifact{
        {Path: "buildfile", Content: []byte("...")},
        {Path: "src/main.ext", Content: []byte("...")},
    }
}

// ... implement other methods
```

## 10. Artifact Model

```go
type Artifact struct {
    Path    string  // Relative path dalam output directory
    Content []byte  // Isi file
}
```

Path digunakan untuk menentukan lokasi output:
- Artefak dengan path yang sama akan menimpa (last-write-wins).
- Artefak dapat menggunakan prefix directory untuk grouping.

## 11. Integration with Pipeline

```
Pipeline.Run()
    ↓
NEIR Model Built
    ↓
adapters.GenerateForNEIR(neir)
    ↓
For each language in GenerationConfig:
    adapter := adapters.Get(lang)
    artifacts = append(artifacts, adapter.Generate*()...)
    ↓
Artifacts written to output directory
```

## 12. Constraints
- Adapter harus thread-safe (register boleh concurrent, generate tidak boleh).
- Artifact Path harus relatif dan menggunakan forward slash.
- Setiap adapter harus menghasilkan minimal: project scaffold, Dockerfile, CI workflow.
- Adapter tidak boleh memiliki dependensi ke adapter lain.

## 13. Future Work
- Plugin-based adapter loading (dynamic registration).
- Adapter versioning dan compatibility checks.
- Shared artifacts (README, LICENSE) across languages.
- Incremental generation (only changed artifacts).

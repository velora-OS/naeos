# NES-008 Registry

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the registry as the metadata catalog used to discover and resolve NAEOS artifacts.

## 3. Scope
The registry covers registration of templates, plugins, blueprints, artifacts, and dependency metadata.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The registry shall support artifact discovery by identifier and category.
- FR-002: The registry shall store compatibility and version metadata.
- FR-003: The registry shall support registration with metadata (version, category, custom metadata).
- FR-004: The registry shall support lookup by name, category, and version.
- FR-005: The registry shall support deregistration of entries.

### 4.2 Non-Functional Requirements
- NFR-001: Registry access shall be auditable.
- NFR-002: Registry operations shall be thread-safe.
- NFR-003: Registry queries shall be performant for typical project sizes.

## 5. Registry Model

### 5.1 Architecture

```
Registry
└── entries: map[string]*Entry
```

Thread-safe via `sync.RWMutex`.

### 5.2 Entry

```go
type Entry struct {
    Name      string
    Version   string
    Category  string
    Component any
    Metadata  map[string]string
}
```

Field | Tipe | Deskripsi
------|------|----------
Name | string | Nama unik entry
Version | string | Versi komponen
Category | string | Kategori (template, plugin, blueprint, service, module)
Component | any | Komponen yang didaftarkan
Metadata | map[string]string | Metadata tambahan

### 5.3 Constructor

```go
func NewRegistry() *Registry
```

## 6. Operations

### 6.1 Registration

```go
func (r *Registry) Register(name string, component any) error
func (r *Registry) RegisterWithMeta(name, version, category string, component any, metadata map[string]string) error
```

- `Register` — pendaftaran sederhana tanpa metadata.
- `RegisterWithMeta` — pendaftaran lengkap dengan versi, kategori, dan metadata.
- Error jika nama kosong atau sudah terdaftar.

### 6.2 Resolution

```go
func (r *Registry) Resolve(name string) (any, error)
func (r *Registry) GetEntry(name string) (*Entry, error)
```

- `Resolve` — mengembalikan komponen berdasarkan nama.
- `GetEntry` — mengembalikan entry lengkap berdasarkan nama.

### 6.3 Deregistration

```go
func (r *Registry) Unregister(name string) error
```

Menghapus entry dari registry. Error jika entry tidak ditemukan.

### 6.4 Query

```go
func (r *Registry) RegisteredEntries() []string
func (r *Registry) FindByCategory(category string) []*Entry
func (r *Registry) FindByVersion(version string) []*Entry
func (r *Registry) Count() int
func (r *Registry) Contains(name string) bool
```

## 7. Usage Example

```go
reg := registry.NewRegistry()

// Register component
reg.RegisterWithMeta("renderer-go", "1.0.0", "template", goRenderer, nil)

// Resolve
component, err := reg.Resolve("renderer-go")

// Query by category
templates := reg.FindByCategory("template")
```

## 8. Acceptance Criteria
- A component can discover a registered artifact without manual lookup.
- Registry metadata is sufficient for version-aware integration.
- Thread-safe operations under concurrent access.

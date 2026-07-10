# NES-036 Template Renderer

## 1. Status
- Status: Draft
- Version: 0.1
- Owner: NAEOS Core Team

## 2. Purpose
Dokumentasi referensi untuk package `internal/generation/renderers` — engine rendering template berbasis `text/template` Go untuk menghasilkan artefak kode dan dokumen.

## 3. Scope
Dokumen ini mencakup interface renderer, implementasi default, fungsi helper, dan tipe data template.

## 4. Normative References
- NES-007 Generator — transformasi desain menjadi artefak
- NES-026 Pipeline — orkestrasi pipeline

## 5. Interface

### 5.1 Renderer

```go
type Renderer interface {
    Render(tmpl string, data any) ([]byte, error)
    RenderNamed(name, tmpl string, data any) ([]byte, error)
}
```

Metode | Deskripsi
-------|----------
Render(tmpl, data) | Me-render template unnamed dengan data
RenderNamed(name, tmpl, data) | Me-render template bernama dengan data

## 6. Implementasi

### 6.1 DefaultRenderer

```go
type DefaultRenderer struct{}
```

Implementasi default dari interface Renderer. Stateless dan thread-safe.

### 6.2 Constructor

```go
func NewRenderer() Renderer
```

Membuat instance baru dari DefaultRenderer.

## 7. Functions

### 7.1 RenderTemplate

```go
func RenderTemplate(name, tmpl string, data any) ([]byte, error)
```

Me-render template string dengan nama dan data tertentu.

Langkah:
1. Parse template dengan nama.
2. Execute template dengan data ke buffer.
3. Mengembalikan hasil sebagai byte slice.

### 7.2 RenderWithFuncs

```go
func RenderWithFuncs(name, tmpl string, data any, funcs template.FuncMap) ([]byte, error)
```

Me-render template dengan custom function map.

```go
funcs := template.FuncMap{
    "toUpper": strings.ToUpper,
    "indent": func(s string) string {
        return "  " + s
    },
}

result, err := renderers.RenderWithFuncs("config", tmpl, data, funcs)
```

## 8. Types

### 8.1 TemplateData

```go
type TemplateData struct {
    Project    string
    Module     string
    Service    string
    Port       int
    Kind       string
    Version    string
    Package    string
    Attributes map[string]string
}
```

Tipe data standar yang digunakan sebagai input template.

Field | Tipe | Deskripsi
------|------|----------
Project | string | Nama proyek
Module | string | Nama modul
Service | string | Nama service
Port | int | Port service
Kind | string | Jenis artefak (http, grpc, dll)
Version | string | Versi artefak
Package | string | Nama package Go
Attributes | map[string]string | Atribut tambahan

## 9. Usage Example

```go
renderer := renderers.NewRenderer()

data := renderers.TemplateData{
    Project: "my-project",
    Module:  "user",
    Package: "user",
}

result, err := renderer.Render(
    "package {{.Project}}/internal/{{.Module}}",
    data,
)
```

## 10. Constraints

- Template harus valid Go text/template syntax.
- Error dari parse dan execute dikembalikan sebagai wrapped error.
- Tidak ada caching template (diparsing ulang setiap kali).
- Tidak ada sandboxing — template memiliki akses penuh ke data.

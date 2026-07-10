# NES-017 Studio

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the user interface and experience layer for interacting with NAEOS projects.

## 3. Scope
The studio covers CLI interface, interactive mode, visualization, and developer experience.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The studio shall provide a CLI interface for all NAEOS operations.
- FR-002: The studio shall support multiple output formats (text, JSON, YAML).
- FR-003: The studio shall provide inspection and preview capabilities.
- FR-004: The studio shall support diagnostic operations (doctor, repair).

### 4.2 Non-Functional Requirements
- NFR-001: The studio shall provide clear and actionable error messages.
- NFR-002: The studio shall be responsive for typical project sizes.

## 5. CLI Commands

| Command | Deskripsi |
|---------|-----------|
| `init` | Inisialisasi proyek baru |
| `run` | Jalankan pipeline lengkap |
| `validate` | Validasi spesifikasi |
| `inspect` | Inspeksi model NEIR |
| `doctor` | Diagnosa lingkungan |
| `repair` | Perbaiki artefak |
| `scaffold` | Buat struktur proyek |
| `export` | Ekspor model |
| `preview` | Pratinjau output |
| `kernel` | Kelola kernel services |
| `version` | Tampilkan versi |

### 5.1 Kernel Subcommands

| Subcommand | Deskripsi |
|------------|-----------|
| `kernel services` | Daftar registered services |
| `kernel metrics` | Tampilkan metrik kernel |
| `kernel events` | Tampilkan event log |
| `kernel publish` | Publish event |
| `kernel subscribe` | Subscribe ke topik |

### 5.2 Output Formats

```bash
# Text (default)
naeos inspect --format text

# JSON
naeos inspect --format json

# YAML
naeos inspect --format yaml
```

## 6. Developer Experience

### 6.1 Quick Start

```bash
# Initialize project
naeos init

# Run pipeline
naeos run specification.yaml

# Validate
naeos validate specification.yaml

# Inspect
naeos inspect
```

### 6.2 Diagnostic Workflow

```bash
# Check environment
naeos doctor

# Repair artifacts
naeos repair
```

## 7. Workflow
1. User invokes a CLI command.
2. CLI parses arguments and flags.
3. CLI dispatches to the appropriate handler.
4. Handler executes the operation.
5. Results are formatted and displayed.

## 8. Acceptance Criteria
- All NAEOS operations are accessible via CLI.
- Output formats (text, JSON, YAML) work correctly.
- Error messages are clear and actionable.
- Diagnostic tools help identify and fix issues.

# NES-010 Knowledge

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the knowledge model used to preserve design context, rationale, and operational history.

## 3. Scope
The knowledge model covers technical decisions, architecture records, policy context, implementation notes, and historical artifacts. Implementation details are documented in NES-037.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: The platform shall store design rationale and decision records.
- FR-002: The platform shall support retrieval of knowledge by topic, component, and version.
- FR-003: The platform shall support 14 node types for comprehensive knowledge capture.
- FR-004: The platform shall support 13 edge types for relationship modeling.
- FR-005: The platform shall support path existence checking between nodes.

### 4.2 Non-Functional Requirements
- NFR-001: Knowledge shall be versioned and searchable.
- NFR-002: Knowledge shall retain provenance for auditability.
- NFR-003: Knowledge graph operations shall be thread-safe.

## 5. Knowledge Forms

### 5.1 Node Types

| Tipe | Deskripsi | Contoh |
|------|-----------|--------|
| decision | Keputusan engineering | "Menggunakan PostgreSQL" |
| requirement | Kebutuhan fungsional | "Harus mendukung 1000 user" |
| rationale | Alasan keputusan | "ACID compliance" |
| component | Komponen sistem | "User Service" |
| policy | Aturan governance | "Health check wajib" |
| implementation | Detail implementasi | "GORM untuk ORM" |
| historical | Riwayat perubahan | "Migrasi DB" |
| service | Service proyek | "auth-service" |
| module | Modul proyek | "user-module" |
| api | Endpoint API | "GET /api/users" |
| storage | Penyimpanan data | "PostgreSQL table" |
| deployment | Konfigurasi deploy | "K8s deployment" |
| testing | Strategi testing | "Unit test > 80%" |
| security | Kontrol keamanan | "JWT auth" |

### 5.2 Edge Types

| Tipe | Deskripsi |
|------|-----------|
| depends_on | Dependensi |
| implements | Implementasi |
| related_to | Relasi umum |
| supersedes | Menggantikan |
| conflicts_with | Konflik |
| contains | Mengandung |
| exposes | Mengekspos |
| connects_to | Koneksi |
| deploys_to | Deploy target |
| tests | Pengujian |
| secures | Keamanan |
| uses | Menggunakan |
| extends | Meng-extend |

## 6. Workflow
1. **Record** the decision or knowledge item as a node.
2. **Connect** related nodes with typed edges.
3. **Attach** provenance and relevant context.
4. **Query** knowledge by topic, component, version, or type.
5. **Trace** paths between related knowledge items.

## 7. Provenance

Setiap knowledge entry harus memiliki provenance yang dapat ditelusuri. Lihat NES-037 untuk detail implementasi provenance tracking.

## 8. Acceptance Criteria
- Key decisions can be retrieved and reviewed at a later time.
- Knowledge entries remain traceable to the relevant artifact or system version.
- Knowledge graph supports queries by topic, component, version, and type.
- Path existence can be verified between any two nodes.

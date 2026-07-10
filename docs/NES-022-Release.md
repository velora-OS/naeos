# NES-022 Release

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the publication and rollout process for NAEOS-generated projects.

## 3. Scope
The release layer covers versioning, changelog generation, release automation, and distribution.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: NAEOS shall support semantic versioning for released artifacts.
- FR-002: NAEOS shall generate changelog from specification changes.
- FR-003: Release artifacts shall be traceable to source specifications.
- FR-004: Release process shall be automatable via CI/CD.

### 4.2 Non-Functional Requirements
- NFR-001: Release artifacts shall be immutable once published.
- NFR-002: Release process shall be auditable and reversible.

## 5. Release Model

### 5.1 Versioning

Menggunakan Semantic Versioning (SemVer):

```
MAJOR.MINOR.PATCH

MAJOR — breaking changes
MINOR — new features (backward compatible)
PATCH — bug fixes
```

### 5.2 Release Artifacts

| Artifact | Deskripsi |
|----------|-----------|
| Source Code | Kode yang di-tag dengan versi |
| Binary | Compiled binaries (jika applicable) |
| Container Image | Docker image dengan versi tag |
| Documentation | Dokumentasi yang sesuai dengan versi |
| Changelog | Daftar perubahan |

### 5.3 Release Workflow

```
Specification Updated
    ↓
NEIR Model Rebuilt
    ↓
Artifacts Regenerated
    ↓
Validation Passed
    ↓
Review Approved
    ↓
Version Bumped
    ↓
Changelog Generated
    ↓
Artifacts Published
    ↓
Release Tagged
```

### 5.4 Changelog Format

```markdown
## [1.2.0] - 2026-07-10

### Added
- New user authentication module
- JWT token validation

### Changed
- Updated database schema for user table

### Fixed
- Fixed race condition in event handler
```

## 6. Release Strategies

### 6.1 Manual Release

```bash
naeos run specification.yaml
naeos validate specification.yaml
# Review and approve
git tag v1.2.0
git push origin v1.2.0
```

### 6.2 Automated Release

```yaml
# .github/workflows/release.yml
on:
  push:
    tags: ['v*']
jobs:
  release:
    steps:
      - run: naeos run specification.yaml
      - run: naeos validate specification.yaml
      - run: docker build -t myapp:${{ github.ref_name }} .
      - run: docker push myapp:${{ github.ref_name }}
```

## 7. Acceptance Criteria
- Released artifacts are versioned according to SemVer.
- Changelog accurately reflects specification changes.
- Release artifacts are traceable to source specifications.
- Release process is automatable via CI/CD.

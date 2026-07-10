# NES-021 Testing

## 1. Status
- Status: Draft
- Version: 0.2
- Owner: NAEOS Core Team

## 2. Purpose
This specification defines the testing quality gates and regression coverage requirements for NAEOS projects.

## 3. Scope
The testing layer covers unit testing, integration testing, validation testing, and test generation.

## 4. Requirements
### 4.1 Functional Requirements
- FR-001: NAEOS shall generate test files for Go modules.
- FR-002: NAEOS shall enforce minimum test coverage requirements.
- FR-003: Tests shall be generated from specification patterns.
- FR-004: Test results shall be auditable through governance.

### 4.2 Non-Functional Requirements
- NFR-001: Generated tests shall be runnable without modification.
- NFR-002: Test coverage shall meet configurable thresholds.

## 5. Testing Model

### 5.1 Test Generation

Generator menghasilkan test files untuk setiap module:

| File | Deskripsi |
|------|-----------|
| `*_test.go` | Unit tests untuk Go packages |
| `handler_test.go` | HTTP handler tests |
| `service_test.go` | Service layer tests |
| `repository_test.go` | Repository tests |

### 5.2 Test Patterns

#### Handler Test
```go
func TestGetUser(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users/1", nil)
    rr := httptest.NewRecorder()

    handler := NewUserHandler(service)
    handler.GetUser(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
}
```

#### Service Test
```go
func TestCreateUser(t *testing.T) {
    repo := &MockUserRepository{}
    service := NewUserService(repo)

    user, err := service.CreateUser(context.Background(), CreateUserRequest{
        Name: "Test User",
    })

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("expected name 'Test User', got '%s'", user.Name)
    }
}
```

### 5.3 Coverage Requirements

| Metric | Minimum |
|--------|---------|
| Unit Test Coverage | > 70% |
| Handler Test Coverage | > 80% |
| Service Test Coverage | > 75% |

## 6. Validation Testing

NEOS validator melakukan validasi berikut pada artefak:

| Check | File Type | Rule |
|-------|-----------|------|
| Content not empty | All | Content harus ada |
| Package declaration | `.go` | Harus ada `package` statement |
| No TODO | `.go` | Tidak ada `// TODO` |
| No placeholder | `.go` | Tidak ada placeholder text |

## 7. Workflow
1. Developer mendefinisikan testing requirements dalam spesifikasi.
2. Generator menghasilkan test files.
3. Validator memeriksa kualitas test artifacts.
4. CI/CD menjalankan tests dan melaporkan coverage.
5. Governance mengevaluasi test results.

## 8. Acceptance Criteria
- Generated tests are runnable without modification.
- Test coverage meets configurable thresholds.
- Validation catches missing tests and poor coverage.
- Test results are auditable through governance.

package parser

import "testing"
func TestNewParserParsesJSON(t *testing.T) {
	p := NewParser()
	input := `{"project":{"name":"demo"},"version":1}`

	doc, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if doc.Data == nil {
		t.Fatal("expected parsed data")
	}
}

func TestNewParserParsesYAML(t *testing.T) {
	p := NewParser()
	input := "project:\n  name: demo\nversion: 1\n"

	doc, err := p.Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
	if doc.Data == nil {
		t.Fatal("expected parsed data")
	}
}

func TestNewParserRejectsInvalidInput(t *testing.T) {
	p := NewParser()
	input := "{unclosed"

	doc, err := p.Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
	if doc != nil {
		t.Fatal("expected nil document on invalid input")
	}
}

func TestParserParsesStructuredSpec(t *testing.T) {
	input := `
project: acme-api
version: 1.0
modules:
  - name: auth
    path: ./internal/auth
services:
  - name: gateway
    kind: http
    port: 8080
`

	doc, err := NewParser().Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc.Project != "acme-api" {
		t.Fatalf("expected project name acme-api, got %q", doc.Project)
	}
	if len(doc.Modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(doc.Modules))
	}
	if doc.Modules[0].Name != "auth" {
		t.Fatalf("expected module name auth, got %q", doc.Modules[0].Name)
	}
	if len(doc.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(doc.Services))
	}
	if doc.Services[0].Name != "gateway" {
		t.Fatalf("expected service name gateway, got %q", doc.Services[0].Name)
	}
}

func TestParserParsesFullSpec(t *testing.T) {
	input := `
project: acme-api
version: 1.0
modules:
  - name: auth
    path: ./internal/auth
    description: Authentication module
    dependencies:
      - crypto
      - storage
  - name: users
    path: ./internal/users
    description: User management
services:
  - name: gateway
    kind: http
    port: 8080
    description: API gateway
    endpoints:
      - method: GET
        path: /api/v1/health
        action: healthCheck
      - method: POST
        path: /api/v1/auth/login
        action: login
architecture:
  pattern: hexagonal
  description: Clean architecture with ports and adapters
  principles:
    - separation of concerns
    - dependency inversion
deployment:
  strategy: rolling
  environments:
    - development
    - staging
    - production
testing:
  strategy: unit-integration
  coverage: 80%
`

	doc, err := NewParser().Parse(input)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if doc.Project != "acme-api" {
		t.Fatalf("expected project name acme-api, got %q", doc.Project)
	}
	if len(doc.Modules) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(doc.Modules))
	}
	if doc.Modules[0].Description != "Authentication module" {
		t.Fatalf("expected module description, got %q", doc.Modules[0].Description)
	}
	if len(doc.Modules[0].Dependencies) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(doc.Modules[0].Dependencies))
	}
	if len(doc.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(doc.Services))
	}
	if doc.Services[0].Description != "API gateway" {
		t.Fatalf("expected service description, got %q", doc.Services[0].Description)
	}
	if len(doc.Services[0].Endpoints) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(doc.Services[0].Endpoints))
	}
	if doc.Services[0].Endpoints[0].Method != "GET" {
		t.Fatalf("expected endpoint method GET, got %q", doc.Services[0].Endpoints[0].Method)
	}
	if doc.Services[0].Endpoints[0].Path != "/api/v1/health" {
		t.Fatalf("expected endpoint path /api/v1/health, got %q", doc.Services[0].Endpoints[0].Path)
	}
	if doc.Architecture == nil {
		t.Fatal("expected architecture to be parsed")
	}
	if doc.Architecture.Pattern != "hexagonal" {
		t.Fatalf("expected architecture pattern hexagonal, got %q", doc.Architecture.Pattern)
	}
	if len(doc.Architecture.Principles) != 2 {
		t.Fatalf("expected 2 principles, got %d", len(doc.Architecture.Principles))
	}
	if doc.Deployment == nil {
		t.Fatal("expected deployment to be parsed")
	}
	if doc.Deployment.Strategy != "rolling" {
		t.Fatalf("expected deployment strategy rolling, got %q", doc.Deployment.Strategy)
	}
	if len(doc.Deployment.Environments) != 3 {
		t.Fatalf("expected 3 environments, got %d", len(doc.Deployment.Environments))
	}
	if doc.Testing == nil {
		t.Fatal("expected testing to be parsed")
	}
	if doc.Testing.Strategy != "unit-integration" {
		t.Fatalf("expected testing strategy unit-integration, got %q", doc.Testing.Strategy)
	}
	if doc.Testing.Coverage != "80%" {
		t.Fatalf("expected testing coverage 80%%, got %q", doc.Testing.Coverage)
	}
}

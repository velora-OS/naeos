package builder

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/specification/resolver"
)

func TestBuilderCreatesNEIRFromResolvedSpec(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{{"name": "auth", "path": "./internal/auth"}},
		"services": []map[string]any{{"name": "gateway", "kind": "http", "port": 8080}},
	}}

	builder := NewBuilder()
	neir, err := builder.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Project == nil || neir.Project.Name != "acme-api" {
		t.Fatalf("expected project acme-api, got %v", neir.Project)
	}
	if len(neir.Modules) != 1 {
		t.Fatalf("expected one module, got %d", len(neir.Modules))
	}
	if len(neir.Services) != 1 {
		t.Fatalf("expected one service, got %d", len(neir.Services))
	}
	if neir.Services[0].Name != "gateway" {
		t.Fatalf("expected service gateway, got %s", neir.Services[0].Name)
	}
	if neir.Services[0].Port != 8080 {
		t.Fatalf("expected service port 8080, got %d", neir.Services[0].Port)
	}
	if neir.Architecture != nil {
		t.Fatalf("expected nil architecture, got %v", neir.Architecture)
	}
}

func TestBuilderExtractsArchitecture(t *testing.T) {
	resolved := &resolver.ResolvedSpec{Context: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{{"name": "core", "path": "./internal/core"}},
		"architecture": map[string]any{"pattern": "clean", "description": "Clean architecture"},
	}}

	b := NewBuilder()
	neir, err := b.Build(resolved)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	if neir.Architecture == nil {
		t.Fatal("expected architecture to be set")
	}
	if neir.Architecture.Pattern != "clean" {
		t.Fatalf("expected pattern clean, got %s", neir.Architecture.Pattern)
	}
	if neir.Architecture.Description != "Clean architecture" {
		t.Fatalf("expected description, got %s", neir.Architecture.Description)
	}
}

func TestBuilderWithNilInput(t *testing.T) {
	b := NewBuilder()
	_, err := b.Build(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

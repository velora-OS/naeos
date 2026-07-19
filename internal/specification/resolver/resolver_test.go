package resolver

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/specification/normalizer"
)

func TestResolverBuildsContextFromNormalizedSpec(t *testing.T) {
	norm := &normalizer.NormalizedSpec{Values: map[string]any{
		"project":  "acme-api",
		"modules":  []map[string]any{{"name": "auth", "path": "./internal/auth"}},
		"services": []map[string]any{{"name": "gateway", "kind": "http", "port": 8080}},
	}}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(norm)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if resolved.Context["project"] != "acme-api" {
		t.Fatalf("expected project acme-api, got %v", resolved.Context["project"])
	}
	if len(resolved.Context["modules"].([]map[string]any)) != 1 {
		t.Fatalf("expected one module, got %d", len(resolved.Context["modules"].([]map[string]any)))
	}
}

func TestResolverResolvesModuleDependencies(t *testing.T) {
	norm := &normalizer.NormalizedSpec{Values: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{
			{"name": "auth", "path": "./internal/auth"},
			{"name": "user", "path": "./internal/user", "dependencies": []any{"auth", "nonexistent"}},
		},
	}}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(norm)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	modules := resolved.Context["modules"].([]map[string]any)
	userMod := modules[1]
	deps := userMod["dependencies"].([]any)
	if len(deps) != 1 {
		t.Fatalf("expected 1 valid dependency, got %d", len(deps))
	}
	if deps[0] != "auth" {
		t.Fatalf("expected dependency 'auth', got %v", deps[0])
	}
}

func TestResolverNormalizesEndpoints(t *testing.T) {
	norm := &normalizer.NormalizedSpec{Values: map[string]any{
		"project": "acme-api",
		"services": []map[string]any{
			{
				"name": "api",
				"kind": "http",
				"port": 8080,
				"endpoints": []map[string]any{
					{"method": "GET", "path": "users"},
				},
			},
		},
	}}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(norm)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	services := resolved.Context["services"].([]map[string]any)
	eps := services[0]["endpoints"].([]map[string]any)
	if eps[0]["path"] != "/users" {
		t.Fatalf("expected path '/users', got %v", eps[0]["path"])
	}
}

func TestResolverPopulatesDefaults(t *testing.T) {
	norm := &normalizer.NormalizedSpec{Values: map[string]any{
		"project": "acme-api",
		"modules": []map[string]any{
			{"name": "auth"},
		},
	}}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(norm)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	if resolved.Context["architecture"] == nil {
		t.Fatal("expected default architecture to be populated")
	}

	modules := resolved.Context["modules"].([]map[string]any)
	if modules[0]["path"] != "./internal/auth" {
		t.Fatalf("expected default path './internal/auth', got %v", modules[0]["path"])
	}
}

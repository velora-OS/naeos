package normalizer

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/specification/parser"
)

func TestNormalizerConvertsParsedSpecToStructuredValues(t *testing.T) {
	doc := &parser.SpecDocument{
		Project:  "acme-api",
		Modules:  []parser.Module{{Name: "auth", Path: "./internal/auth"}},
		Services: []parser.Service{{Name: "gateway", Kind: "http", Port: 8080}},
	}

	normalizer := NewNormalizer()
	normalized, err := normalizer.Normalize(doc)
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}
	if normalized.Values["project"] != "acme-api" {
		t.Fatalf("expected project acme-api, got %v", normalized.Values["project"])
	}
	if len(normalized.Values["modules"].([]map[string]any)) != 1 {
		t.Fatalf("expected one module, got %d", len(normalized.Values["modules"].([]map[string]any)))
	}
	if normalized.Values["services"].([]map[string]any)[0]["name"] != "gateway" {
		t.Fatalf("expected gateway service, got %v", normalized.Values["services"].([]map[string]any)[0]["name"])
	}
}

func TestNormalizerHandlesArchitectureDeploymentTesting(t *testing.T) {
	doc := &parser.SpecDocument{
		Project: "full-api",
		Modules: []parser.Module{
			{Name: "auth", Path: "./internal/auth", Description: "Auth module", Dependencies: []string{"crypto"}},
		},
		Services: []parser.Service{
			{Name: "gateway", Kind: "http", Port: 8080, Description: "API gateway", Endpoints: []parser.Endpoint{
				{Method: "GET", Path: "/health", Action: "healthCheck"},
			}},
		},
		Architecture: &parser.Architecture{
			Pattern:     "hexagonal",
			Description: "Clean architecture",
			Principles:  []string{"separation of concerns", "dependency inversion"},
		},
		Deployment: &parser.Deployment{
			Strategy:     "rolling",
			Environments: []string{"staging", "production"},
		},
		Testing: &parser.Testing{
			Strategy: "unit-integration",
			Coverage: "80%",
		},
	}

	normalized, err := NewNormalizer().Normalize(doc)
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}

	arch := normalized.Values["architecture"].(map[string]any)
	if arch["pattern"] != "hexagonal" {
		t.Fatalf("expected pattern hexagonal, got %v", arch["pattern"])
	}
	principles := arch["principles"].([]string)
	if len(principles) != 2 {
		t.Fatalf("expected 2 principles, got %d", len(principles))
	}

	deploy := normalized.Values["deployment"].(map[string]any)
	if deploy["strategy"] != "rolling" {
		t.Fatalf("expected strategy rolling, got %v", deploy["strategy"])
	}
	envs := deploy["environments"].([]map[string]any)
	if len(envs) != 2 {
		t.Fatalf("expected 2 environments, got %d", len(envs))
	}

	test := normalized.Values["testing"].(map[string]any)
	if test["strategy"] != "unit-integration" {
		t.Fatalf("expected strategy unit-integration, got %v", test["strategy"])
	}
	if test["coverage"] != "80%" {
		t.Fatalf("expected coverage 80%%, got %v", test["coverage"])
	}

	modules := normalized.Values["modules"].([]map[string]any)
	if modules[0]["description"] != "Auth module" {
		t.Fatalf("expected module description, got %v", modules[0]["description"])
	}
	deps := modules[0]["dependencies"].([]string)
	if deps[0] != "crypto" {
		t.Fatalf("expected dependency crypto, got %v", deps[0])
	}

	services := normalized.Values["services"].([]map[string]any)
	if services[0]["description"] != "API gateway" {
		t.Fatalf("expected service description, got %v", services[0]["description"])
	}
	eps := services[0]["endpoints"].([]map[string]any)
	if eps[0]["method"] != "GET" {
		t.Fatalf("expected endpoint method GET, got %v", eps[0]["method"])
	}
}

func TestNormalizerNilDocumentReturnsError(t *testing.T) {
	_, err := NewNormalizer().Normalize(nil)
	if err == nil {
		t.Fatal("expected error for nil document")
	}
}

func TestNormalizerNonSpecDocumentReturnsSource(t *testing.T) {
	normalized, err := NewNormalizer().Normalize("plain string")
	if err != nil {
		t.Fatalf("Normalize returned error: %v", err)
	}
	if normalized.Values["source"] != "plain string" {
		t.Fatalf("expected source to be plain string, got %v", normalized.Values["source"])
	}
}

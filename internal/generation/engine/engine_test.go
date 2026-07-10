package engine

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
)

func TestGeneratorCreatesArtifactsFromNEIR(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "acme-api"},
		Modules: []module.Module{{Name: "auth", Path: "./internal/auth"}},
	}

	engine := NewEngine()
	artifacts, err := engine.Generate(neir)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if len(artifacts) < 2 {
		t.Fatalf("expected at least two artifacts, got %d", len(artifacts))
	}
	if artifacts[0].Path != "README.md" {
		t.Fatalf("expected README artifact first, got %s", artifacts[0].Path)
	}
}

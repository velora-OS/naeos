package validator

import (
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/service"
)

func TestValidatorAcceptsValidNEIR(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "acme-api"},
		Modules: []module.Module{{Name: "auth", Path: "./internal/auth"}},
	}
	v := NewValidator()
	if err := v.Validate(neir); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestValidatorRejectsIncompleteNEIR(t *testing.T) {
	neir := &model.NEIR{Project: &project.Project{}, Modules: []module.Module{}}
	v := NewValidator()
	if err := v.Validate(neir); err == nil {
		t.Fatalf("expected validation error for incomplete NEIR")
	}
}

func TestValidatorRejectsNilInput(t *testing.T) {
	v := NewValidator()
	if err := v.Validate(nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestValidatorRejectsDuplicateModuleNames(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "test"},
		Modules: []module.Module{
			{Name: "auth", Path: "./internal/auth"},
			{Name: "auth", Path: "./internal/auth2"},
		},
	}
	result := ValidateDetailed(neir)
	if result.Valid {
		t.Fatal("expected validation to fail for duplicate module names")
	}
}

func TestValidatorRejectsEmptyModuleName(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "test"},
		Modules: []module.Module{{Name: "", Path: "./internal/empty"}},
	}
	result := ValidateDetailed(neir)
	if result.Valid {
		t.Fatal("expected validation to fail for empty module name")
	}
}

func TestValidatorRejectsEmptyModulePath(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "test"},
		Modules: []module.Module{{Name: "core", Path: ""}},
	}
	result := ValidateDetailed(neir)
	if result.Valid {
		t.Fatal("expected validation to fail for empty module path")
	}
}

func TestValidatorRejectsInvalidServicePort(t *testing.T) {
	neir := &model.NEIR{
		Project:  &project.Project{Name: "test"},
		Modules:  []module.Module{{Name: "core", Path: "./internal/core"}},
		Services: []service.Service{{Name: "api", Port: 99999}},
	}
	result := ValidateDetailed(neir)
	if result.Valid {
		t.Fatal("expected validation to fail for invalid service port")
	}
}

func TestValidatorReportsMultipleErrors(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{},
		Modules: []module.Module{},
	}
	result := ValidateDetailed(neir)
	if result.Valid {
		t.Fatal("expected validation to fail")
	}
	if len(result.Errors) < 2 {
		t.Fatalf("expected multiple errors, got %d", len(result.Errors))
	}
}

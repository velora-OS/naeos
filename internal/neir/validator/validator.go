package validator

import (
	"fmt"
	"strings"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

type Validator interface {
	Validate(neir any) error
}

type ValidationResult struct {
	Valid  bool
	Errors []string
}

type DefaultValidator struct{}

func NewValidator() Validator {
	return DefaultValidator{}
}

func (DefaultValidator) Validate(neir any) error {
	result := ValidateDetailed(neir)
	if !result.Valid {
		return fmt.Errorf("validation failed: %s", strings.Join(result.Errors, "; "))
	}
	return nil
}

func ValidateDetailed(neir any) ValidationResult {
	result := ValidationResult{Valid: true}

	if neir == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "neir is nil")
		return result
	}

	neirStruct, ok := neir.(*model.NEIR)
	if !ok {
		result.Valid = false
		result.Errors = append(result.Errors, "neir is not a model.NEIR")
		return result
	}

	if neirStruct.Project == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "project must be set")
	} else if neirStruct.Project.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "project name must not be empty")
	}

	if len(neirStruct.Modules) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "must contain at least one module")
	}

	for i, mod := range neirStruct.Modules {
		if mod.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module[%d] name must not be empty", i))
		}
		if mod.Path == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module[%d] (%s) path must not be empty", i, mod.Name))
		}
	}

	seenModules := make(map[string]int)
	for i, mod := range neirStruct.Modules {
		if prev, exists := seenModules[mod.Name]; exists {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module name %q is duplicated at index %d and %d", mod.Name, prev, i))
		}
		seenModules[mod.Name] = i
	}

	for i, svc := range neirStruct.Services {
		if svc.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service[%d] name must not be empty", i))
		}
		if svc.Port < 0 || svc.Port > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service[%d] (%s) port must be between 0 and 65535, got %d", i, svc.Name, svc.Port))
		}
	}

	if neirStruct.Metadata != nil {
		if neirStruct.Metadata.NEIRVersion == "" {
			result.Errors = append(result.Errors, "metadata.neir_version is recommended")
		}
	}

	if neirStruct.Generation != nil {
		if len(neirStruct.Generation.Languages) == 0 {
			result.Errors = append(result.Errors, "generation.languages should not be empty when generation is set")
		}
		for _, lang := range neirStruct.Generation.Languages {
			if !language.IsValid(lang) {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("generation.language %q is not a supported language (valid: go, typescript, python, java, rust)", lang))
			}
		}
	}

	return result
}

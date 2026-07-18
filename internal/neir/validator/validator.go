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
	Warns  []string
}

type DefaultValidator struct{}

func NewValidator() Validator {
	return DefaultValidator{}
}

func (DefaultValidator) Validate(neir any) error {
	result := ValidateDetailed(neir)
	if !result.Valid {
		return fmt.Errorf("validation failed:\n  - %s", strings.Join(result.Errors, "\n  - "))
	}
	return nil
}

func ValidateDetailed(neir any) ValidationResult {
	result := ValidationResult{Valid: true}

	if neir == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "NEIR model is nil — ensure the specification was parsed correctly")
		return result
	}

	neirStruct, ok := neir.(*model.NEIR)
	if !ok {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("expected *model.NEIR, got %T — this is an internal error", neir))
		return result
	}

	if neirStruct.Project == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "project is required — add a 'project:' field to your specification")
	} else if neirStruct.Project.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "project name must not be empty — set 'project: <name>' in your specification")
	}

	if len(neirStruct.Modules) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, "at least one module is required — add a 'modules:' section to your specification")
	}

	for i, mod := range neirStruct.Modules {
		if mod.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module[%d] name is required — each module needs a 'name:' field", i))
		}
		if mod.Path == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("module %q (index %d) path is required — add a 'path:' field (e.g., ./internal/%s)", mod.Name, i, mod.Name))
		}
	}

	seenModules := make(map[string]int)
	for i, mod := range neirStruct.Modules {
		if prev, exists := seenModules[mod.Name]; exists {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("duplicate module name %q at index %d and %d — module names must be unique", mod.Name, prev, i))
		}
		seenModules[mod.Name] = i
	}

	for i, svc := range neirStruct.Services {
		if svc.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service[%d] name is required — each service needs a 'name:' field", i))
		}
		if svc.Port < 0 || svc.Port > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("service %q port %d is out of range — must be between 0 and 65535", svc.Name, svc.Port))
		}
	}

	seenPorts := make(map[int]string)
	for _, svc := range neirStruct.Services {
		if svc.Port > 0 {
			if prev, exists := seenPorts[svc.Port]; exists {
				result.Warns = append(result.Warns, fmt.Sprintf("service %q and %q share port %d — this may cause conflicts", prev, svc.Name, svc.Port))
			}
			seenPorts[svc.Port] = svc.Name
		}
	}

	// Validate Architecture
	if neirStruct.Architecture != nil {
		if neirStruct.Architecture.Pattern == "" {
			result.Valid = false
			result.Errors = append(result.Errors, "architecture.pattern is required when architecture section is present")
		} else {
			validPatterns := map[string]bool{
				"layered":    true,
				"clean":      true,
				"hexagonal":  true,
				"microkernel":true,
				"event-driven":true,
				"cqrs":       true,
				"monolith":   true,
			}
			if !validPatterns[neirStruct.Architecture.Pattern] {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported architecture pattern %q — supported: layered, clean, hexagonal, microkernel, event-driven, cqrs, monolith", neirStruct.Architecture.Pattern))
			}
		}
	}

	// Validate Deployment
	if neirStruct.Deployment != nil {
		if neirStruct.Deployment.Strategy == "" {
			result.Valid = false
			result.Errors = append(result.Errors, "deployment.strategy is required when deployment section is present")
		} else {
			validStrategies := map[string]bool{
				"rolling":   true,
				"blue-green":true,
				"canary":    true,
				"recreate":  true,
			}
			if !validStrategies[neirStruct.Deployment.Strategy] {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported deployment strategy %q — supported: rolling, blue-green, canary, recreate", neirStruct.Deployment.Strategy))
			}
		}

		if len(neirStruct.Deployment.Environments) == 0 {
			result.Warns = append(result.Warns, "deployment.environments is empty — consider specifying target environments")
		}
	}

	// Validate Testing
	if neirStruct.Testing != nil {
		if neirStruct.Testing.Strategy == "" {
			result.Valid = false
			result.Errors = append(result.Errors, "testing.strategy is required when testing section is present")
		} else {
			validStrategies := map[string]bool{
				"unit":      true,
				"integration":true,
				"e2e":       true,
				"contract":  true,
			}
			if !validStrategies[neirStruct.Testing.Strategy] {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported testing strategy %q — supported: unit, integration, e2e, contract", neirStruct.Testing.Strategy))
			}
		}

		if neirStruct.Testing.Coverage != "" {
			// Basic validation for coverage percentage format
			if len(neirStruct.Testing.Coverage) > 0 && neirStruct.Testing.Coverage[len(neirStruct.Testing.Coverage)-1] != '%' {
				result.Warns = append(result.Warns, "testing.coverage should be a percentage value (e.g., '80%')")
			}
		}
	}

	// Validate Cloud
	if neirStruct.Cloud != nil {
		if neirStruct.Cloud.Provider == "" {
			result.Valid = false
			result.Errors = append(result.Errors, "cloud.provider is required when cloud section is present")
		} else {
			validProviders := map[string]bool{
				"aws":       true,
				"gcp":       true,
				"azure":     true,
				"digitalocean": true,
			}
			if !validProviders[neirStruct.Cloud.Provider] {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported cloud provider %q — supported: aws, gcp, azure, digitalocean", neirStruct.Cloud.Provider))
			}
		}

		if neirStruct.Cloud.Region == "" {
			result.Warns = append(result.Warns, "cloud.region is not specified — using provider's default region")
		}

		for i, svc := range neirStruct.Cloud.Services {
			if svc.Name == "" {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("cloud.services[%d] name is required", i))
			}
			if svc.Type == "" {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("cloud.services[%d] type is required", i))
			} else {
				validTypes := map[string]bool{
					"compute": true,
					"storage": true,
					"database":true,
					"cache":   true,
					"queue":   true,
				}
				if !validTypes[svc.Type] {
					result.Valid = false
					result.Errors = append(result.Errors, fmt.Sprintf("unsupported cloud service type %q for cloud.services[%d] — supported: compute, storage, database, cache, queue", svc.Type, i))
				}
			}
			if svc.Tier != "" {
				validTiers := map[string]bool{
					"small":  true,
					"medium": true,
					"large":  true,
				}
				if !validTiers[svc.Tier] {
					result.Warns = append(result.Warns, fmt.Sprintf("cloud.services[%d] tier %q is not standard — recommended: small, medium, large", i, svc.Tier))
				}
			}
		}

		if neirStruct.Cloud.Scaling != nil {
			if neirStruct.Cloud.Scaling.MinReplicas < 0 {
				result.Valid = false
				result.Errors = append(result.Errors, "cloud.scaling.min_replicas must be non-negative")
			}
			if neirStruct.Cloud.Scaling.MaxReplicas < neirStruct.Cloud.Scaling.MinReplicas {
				result.Valid = false
				result.Errors = append(result.Errors, "cloud.scaling.max_replicas must be greater than or equal to min_replicas")
			}
			if neirStruct.Cloud.Scaling.TargetCPUUtilizationPercent != nil {
				if *neirStruct.Cloud.Scaling.TargetCPUUtilizationPercent < 1 || *neirStruct.Cloud.Scaling.TargetCPUUtilizationPercent > 100 {
					result.Valid = false
					result.Errors = append(result.Errors, "cloud.scaling.target_cpu must be between 1 and 100")
				}
			}
		}
	}

	// Validate Plugins
	for i, plugin := range neirStruct.Plugins {
		if plugin.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("plugins[%d] name is required", i))
		}
		if plugin.Source == "" {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("plugins[%d] source is required", i))
		}
	}

	// Validate AI
	if neirStruct.AI != nil {
		if neirStruct.AI.ContextType != "" {
			validContextTypes := map[string]bool{
				"full":        true,
				"summary":     true,
				"dependencies":true,
			}
			if !validContextTypes[neirStruct.AI.ContextType] {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported AI context_type %q — supported: full, summary, dependencies", neirStruct.AI.ContextType))
			}
		}

		validEnrichments := map[string]bool{
			"security":     true,
			"performance":  true,
			"testing":      true,
			"documentation":true,
			"performance":  true,
		}
		for _, enrichment := range neirStruct.AI.Enrichment {
			if !validEnrichments[enrichment] {
				result.Warns = append(result.Warns, fmt.Sprintf("AI enrichment %q is not a standard value — common values: security, performance, testing, documentation", enrichment))
			}
		}
	}

	// Validate Metadata
	if neirStruct.Metadata != nil {
		if neirStruct.Metadata.NEIRVersion == "" {
			result.Warns = append(result.Warns, "metadata.neir_version is recommended for traceability")
		}
		if neirStruct.Metadata.SchemaVersion == "" {
			result.Warns = append(result.Warns, "metadata.schema_version is recommended for schema tracking")
		}
	}

	// Validate Generation
	if neirStruct.Generation != nil {
		if len(neirStruct.Generation.Languages) == 0 {
			result.Warns = append(result.Warns, "generation.languages is empty — defaulting to Go")
		}
		for _, lang := range neirStruct.Generation.Languages {
			if !language.IsValid(lang) {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("unsupported language %q — supported: go, typescript, python, java, rust", lang))
			}
		}
		if neirStruct.Generation.OutputDir == "" {
			result.Warns = append(result.Warns, "generation.output_dir is not specified — using default output directory")
		}
	}

	return result
}

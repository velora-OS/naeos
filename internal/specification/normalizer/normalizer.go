package normalizer

import (
	"fmt"

	"github.com/NAEOS-foundation/naeos/internal/specification/parser"
)

type Normalizer interface {
	Normalize(doc any) (*NormalizedSpec, error)
}

type NormalizedSpec struct {
	Values map[string]any
}

type DefaultNormalizer struct{}

func NewNormalizer() Normalizer {
	return DefaultNormalizer{}
}

func (DefaultNormalizer) Normalize(doc any) (*NormalizedSpec, error) {
	if doc == nil {
		return nil, fmt.Errorf("document is nil")
	}

	specDoc, ok := doc.(*parser.SpecDocument)
	if !ok {
		return &NormalizedSpec{Values: map[string]any{"source": doc}}, nil
	}

	modules := normalizeModules(specDoc.Modules)
	services := normalizeServices(specDoc.Services)

	result := map[string]any{
		"project":  specDoc.Project,
		"modules":  modules,
		"services": services,
		"source":   specDoc,
	}

	if specDoc.Architecture != nil {
		result["architecture"] = normalizeArchitecture(specDoc.Architecture)
	}
	if specDoc.Deployment != nil {
		result["deployment"] = normalizeDeployment(specDoc.Deployment)
	}
	if specDoc.Testing != nil {
		result["testing"] = normalizeTesting(specDoc.Testing)
	}
	if specDoc.Generation != nil {
		result["generation"] = normalizeGeneration(specDoc.Generation)
	}

	return &NormalizedSpec{Values: result}, nil
}

func normalizeModules(modules []parser.Module) []map[string]any {
	result := make([]map[string]any, 0, len(modules))
	for _, m := range modules {
		entry := map[string]any{
			"name": m.Name,
			"path": m.Path,
		}
		if m.Description != "" {
			entry["description"] = m.Description
		}
		if len(m.Dependencies) > 0 {
			entry["dependencies"] = m.Dependencies
		}
		result = append(result, entry)
	}
	return result
}

func normalizeServices(services []parser.Service) []map[string]any {
	result := make([]map[string]any, 0, len(services))
	for _, s := range services {
		entry := map[string]any{
			"name": s.Name,
			"kind": s.Kind,
			"port": s.Port,
		}
		if s.Description != "" {
			entry["description"] = s.Description
		}
		if len(s.Endpoints) > 0 {
			eps := make([]map[string]any, 0, len(s.Endpoints))
			for _, ep := range s.Endpoints {
				eps = append(eps, map[string]any{
					"method": ep.Method,
					"path":   ep.Path,
					"action": ep.Action,
				})
			}
			entry["endpoints"] = eps
		}
		result = append(result, entry)
	}
	return result
}

func normalizeArchitecture(arch *parser.Architecture) map[string]any {
	result := map[string]any{
		"pattern":     arch.Pattern,
		"description": arch.Description,
	}
	if len(arch.Principles) > 0 {
		result["principles"] = arch.Principles
	}
	return result
}

func normalizeDeployment(deploy *parser.Deployment) map[string]any {
	result := map[string]any{
		"strategy": deploy.Strategy,
	}
	if len(deploy.Environments) > 0 {
		envs := make([]map[string]any, 0, len(deploy.Environments))
		for _, env := range deploy.Environments {
			envs = append(envs, map[string]any{"name": env})
		}
		result["environments"] = envs
	}
	return result
}

func normalizeGeneration(gen *parser.Generation) map[string]any {
	result := map[string]any{}
	if len(gen.Languages) > 0 {
		result["languages"] = gen.Languages
	}
	if gen.OutputDir != "" {
		result["output_dir"] = gen.OutputDir
	}
	if gen.ModuleDir != "" {
		result["module_dir"] = gen.ModuleDir
	}
	return result
}

func normalizeTesting(test *parser.Testing) map[string]any {
	result := map[string]any{
		"strategy": test.Strategy,
	}
	if test.Coverage != "" {
		result["coverage"] = test.Coverage
	}
	return result
}

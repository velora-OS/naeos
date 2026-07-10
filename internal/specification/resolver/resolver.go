package resolver

import (
	"fmt"

	"github.com/NAEOS-foundation/naeos/internal/specification/normalizer"
)

type Resolver interface {
	Resolve(spec any) (*ResolvedSpec, error)
}

type ResolvedSpec struct {
	Context map[string]any
}

type DefaultResolver struct{}

func NewResolver() Resolver {
	return DefaultResolver{}
}

func (DefaultResolver) Resolve(spec any) (*ResolvedSpec, error) {
	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	normalized, ok := spec.(*normalizer.NormalizedSpec)
	if !ok {
		return &ResolvedSpec{Context: map[string]any{"resolved": true}}, nil
	}

	context := map[string]any{}
	for key, value := range normalized.Values {
		context[key] = value
	}

	resolveModuleDependencies(context)
	resolveServiceEndpoints(context)
	populateDefaults(context)

	return &ResolvedSpec{Context: context}, nil
}

func resolveModuleDependencies(context map[string]any) {
	rawModules, exists := context["modules"]
	if !exists {
		return
	}

	modules, ok := rawModules.([]map[string]any)
	if !ok {
		return
	}

	moduleNames := make(map[string]bool, len(modules))
	for _, m := range modules {
		if name, ok := m["name"].(string); ok {
			moduleNames[name] = true
		}
	}

	for _, m := range modules {
		if deps, ok := m["dependencies"].([]any); ok {
			validDeps := make([]any, 0, len(deps))
			for _, d := range deps {
				if depName, ok := d.(string); ok {
					if moduleNames[depName] {
						validDeps = append(validDeps, d)
					}
				}
			}
			m["dependencies"] = validDeps
		}
	}
}

func resolveServiceEndpoints(context map[string]any) {
	rawServices, exists := context["services"]
	if !exists {
		return
	}

	services, ok := rawServices.([]map[string]any)
	if !ok {
		return
	}

	for _, svc := range services {
		rawEndpoints, ok := svc["endpoints"]
		if !ok {
			continue
		}
		endpoints, ok := rawEndpoints.([]map[string]any)
		if !ok {
			continue
		}
		for _, ep := range endpoints {
			if method, ok := ep["method"].(string); ok {
				ep["method"] = fmt.Sprintf("%s", method)
			}
			if path, ok := ep["path"].(string); ok && path != "" {
				if path[0] != '/' {
					ep["path"] = "/" + path
				}
			}
		}
	}
}

func populateDefaults(context map[string]any) {
	if _, exists := context["architecture"]; !exists {
		context["architecture"] = map[string]any{
			"pattern":     "layered",
			"description": "default layered architecture",
		}
	}

	rawModules, exists := context["modules"]
	if exists {
		if modules, ok := rawModules.([]map[string]any); ok {
			for _, m := range modules {
				if _, hasPath := m["path"]; !hasPath {
					if name, ok := m["name"].(string); ok {
						m["path"] = fmt.Sprintf("./internal/%s", name)
					}
				}
				if _, hasDeps := m["dependencies"]; !hasDeps {
					m["dependencies"] = []any{}
				}
			}
		}
	}
}

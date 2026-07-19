package adapters

import (
	"fmt"
	"strings"
	"time"

	"github.com/NAEOS-foundation/naeos/internal/compiler"
	"github.com/NAEOS-foundation/naeos/internal/neir/model"
)

type geminiAdapter struct{}

func NewGeminiAdapter() compiler.Adapter {
	return &geminiAdapter{}
}

func (a *geminiAdapter) Target() compiler.Target {
	return compiler.TargetGemini
}

func (a *geminiAdapter) Compile(neir *model.NEIR) (*compiler.CompiledOutput, error) {
	if neir == nil {
		return nil, fmt.Errorf("nil NEIR")
	}

	var files []compiler.OutputFile

	geminiFile := a.buildGeminiConfig(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".gemini/CONFIG.md",
		Content: geminiFile,
		Kind:    "instructions",
	})

	contextFile := a.buildContextFile(neir)
	files = append(files, compiler.OutputFile{
		Path:    ".gemini/context.md",
		Content: contextFile,
		Kind:    "context",
	})

	projectName := "unknown"
	if neir.Project != nil {
		projectName = neir.Project.Name
	}

	return &compiler.CompiledOutput{
		Target:     compiler.TargetGemini,
		Files:      files,
		Summary:    fmt.Sprintf("Generated %d files for Gemini CLI (%s)", len(files), projectName),
		CompiledAt: time.Now(),
	}, nil
}

func (a *geminiAdapter) buildGeminiConfig(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Gemini CLI Configuration\n\n")

	if neir.Project != nil {
		fmt.Fprintf(&sb, "Project: %s\n", neir.Project.Name)
		if neir.Project.Description != "" {
			fmt.Fprintf(&sb, "Description: %s\n", neir.Project.Description)
		}
	}

	if neir.Architecture != nil {
		fmt.Fprintf(&sb, "\nArchitecture: %s\n", neir.Architecture.Pattern)
	}

	sb.WriteString("\n## Project Structure\n\n")
	if len(neir.Modules) > 0 {
		for _, m := range neir.Modules {
			fmt.Fprintf(&sb, "- `%s` → `%s`\n", m.Name, m.Path)
			if m.Description != "" {
				fmt.Fprintf(&sb, "  %s\n", m.Description)
			}
		}
	}

	if len(neir.Services) > 0 {
		sb.WriteString("\n## Services\n\n")
		for _, s := range neir.Services {
			fmt.Fprintf(&sb, "- %s (%s, port %d)\n", s.Name, s.Kind, s.Port)
		}
	}

	sb.WriteString("\n## Guidelines\n\n")
	sb.WriteString("- Follow established patterns in the codebase\n")
	sb.WriteString("- Write clean, maintainable code\n")
	sb.WriteString("- Include proper error handling\n")
	sb.WriteString("- Add tests for new functionality\n")

	return sb.String()
}

func (a *geminiAdapter) buildContextFile(neir *model.NEIR) string {
	var sb strings.Builder
	sb.WriteString("# Gemini Context\n\n")

	if len(neir.Components) > 0 {
		sb.WriteString("## Components\n\n")
		for _, c := range neir.Components {
			fmt.Fprintf(&sb, "- %s (%s) in module %s\n", c.Name, c.Kind, c.Module)
		}
	}

	if len(neir.APIs) > 0 {
		sb.WriteString("\n## APIs\n\n")
		for _, api := range neir.APIs {
			fmt.Fprintf(&sb, "### %s (%s)\n", api.Name, api.Protocol)
			for _, ep := range api.Endpoints {
				fmt.Fprintf(&sb, "- %s %s: %s\n", ep.Method, ep.Path, ep.Summary)
			}
		}
	}

	if neir.Testing != nil {
		fmt.Fprintf(&sb, "\n## Testing: %s strategy\n", neir.Testing.Strategy)
		if len(neir.Testing.Frameworks) > 0 {
			fmt.Fprintf(&sb, "Frameworks: %s\n", strings.Join(neir.Testing.Frameworks, ", "))
		}
	}

	return sb.String()
}

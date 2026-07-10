package engine

import (
	"fmt"
	"strings"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

type GeneratorEngine interface {
	Generate(neir any) ([]Artifact, error)
	GenerateForLanguage(neir *model.NEIR, lang language.Language) ([]Artifact, error)
}

type Artifact struct {
	Path    string
	Content []byte
}

type DefaultEngine struct{}

func NewEngine() GeneratorEngine {
	return DefaultEngine{}
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "/", "-")
	value = strings.ReplaceAll(value, "_", "-")
	return value
}

func (DefaultEngine) Generate(neir any) ([]Artifact, error) {
	if neir == nil {
		return nil, fmt.Errorf("neir is nil")
	}

	projectName := fmt.Sprint(neir)
	var modules []any
	var services []any
	var archPattern string
	var deployStrategy string

	if neirMap, ok := neir.(map[string]any); ok {
		if project, ok := neirMap["project"]; ok {
			projectName = fmt.Sprint(project)
		}
		if rawModules, ok := neirMap["modules"].([]any); ok {
			modules = rawModules
		}
		if rawServices, ok := neirMap["services"].([]any); ok {
			services = rawServices
		}
	} else if neirStruct, ok := neir.(*model.NEIR); ok {
		if neirStruct.Project != nil {
			projectName = neirStruct.Project.Name
		}
		for _, m := range neirStruct.Modules {
			modules = append(modules, map[string]any{"name": m.Name, "path": m.Path})
		}
		for _, s := range neirStruct.Services {
			services = append(services, map[string]any{"name": s.Name, "port": s.Port, "kind": string(s.Kind)})
		}
		if neirStruct.Architecture != nil {
			archPattern = string(neirStruct.Architecture.Pattern)
		}
		if neirStruct.Deployment != nil {
			deployStrategy = string(neirStruct.Deployment.Strategy)
		}
	}

	slug := slugify(projectName)

	artifacts := []Artifact{{
		Path: "README.md",
		Content: []byte(fmt.Sprintf("# %s\n\nGenerated from NAEOS pipeline.\n\n## Overview\n\nThis project was scaffolded with NAEOS.\n\n## Project Structure\n\n- cmd/app/main.go - application entrypoint\n- Dockerfile - container build definition\n- .github/workflows/ci.yml - CI workflow\n- spec.yaml - source specification\n\n## Quick Start\n\n1. Review spec.yaml\n2. Run `go test ./...`\n3. Build the app with `go build ./cmd/app`\n4. Run the binary with `./app`\n\n## Deployment\n\nThe generated Dockerfile and CI workflow provide a starting point for shipping the service in a containerized environment.\n", projectName)),
	}, {
		Path:    "Dockerfile",
		Content: []byte("FROM golang:1.22-alpine\nWORKDIR /app\nCOPY . .\nRUN go build ./cmd/app\nCMD [\"/app/app\"]\n"),
	}, {
		Path:    ".github/workflows/ci.yml",
		Content: []byte("name: ci\n\non: [push, pull_request]\n\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-go@v5\n        with:\n          go-version: '1.22'\n      - run: go test ./...\n"),
	}, {
		Path:    "go.mod",
		Content: []byte(fmt.Sprintf("module github.com/example/%s\n\ngo 1.22\n", slug)),
	}, {
		Path:    "cmd/app/main.go",
		Content: []byte(fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from %s\")\n}\n", projectName)),
	}}

	if deployStrategy != "" {
		artifacts = append(artifacts, Artifact{
			Path:    "docker-compose.yml",
			Content: []byte("version: '3.8'\nservices:\n  app:\n    build: .\n    ports:\n      - '8080:8080'\n    deploy:\n      replicas: 2\n"),
		})
	}

	if archPattern != "" {
		artifacts = append(artifacts, Artifact{
			Path:    "docs/architecture.md",
			Content: []byte(fmt.Sprintf("# Architecture\n\nPattern: %s\n\n## Overview\n\nThis document describes the architectural decisions for %s.\n", archPattern, projectName)),
		})
	}

	for _, module := range modules {
		var name string
		var path string

		if moduleMap, ok := module.(map[string]any); ok {
			name = fmt.Sprint(moduleMap["name"])
			path = fmt.Sprint(moduleMap["path"])
		} else if moduleStruct, ok := module.(map[string]string); ok {
			name = moduleStruct["name"]
			path = moduleStruct["path"]
		}

		if name == "" {
			continue
		}
		if path == "" {
			path = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		}
		moduleDir := strings.TrimPrefix(path, "./")
		if moduleDir == "" {
			moduleDir = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		}
		pkg := slugify(name)
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/README.md", moduleDir),
			Content: []byte(fmt.Sprintf("# %s\n\nModule for %s project.\n", name, projectName)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/package.go", moduleDir),
			Content: []byte(fmt.Sprintf("package %s\n\n// %s module placeholder.\n", pkg, name)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/config.yaml", moduleDir),
			Content: []byte(fmt.Sprintf("name: %s\nmodule: %s\n", name, name)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/handler.go", moduleDir),
			Content: []byte(fmt.Sprintf("package %s\n\ntype Handler struct {\n\tservice Service\n}\n\nfunc NewHandler(service Service) *Handler {\n\treturn &Handler{service: service}\n}\n", pkg)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/repository.go", moduleDir),
			Content: []byte(fmt.Sprintf("package %s\n\ntype Repository interface {\n\tList() []string\n}\n", pkg)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/service.go", moduleDir),
			Content: []byte(fmt.Sprintf("package %s\n\ntype Service interface {\n\tHandle() string\n}\n", pkg)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/domain/model.go", moduleDir),
			Content: []byte("package domain\n\ntype Model struct {\n\tName string\n}\n"),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/http/handler.go", moduleDir),
			Content: []byte(fmt.Sprintf("package http\n\nimport \"net/http\"\n\ntype Handler struct{}\n\nfunc (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n\tw.Write([]byte(\"handler for %s\"))\n}\n", name)),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/http/router.go", moduleDir),
			Content: []byte("package http\n\ntype Router struct{}\n"),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/middleware/logging.go", moduleDir),
			Content: []byte("package middleware\n\nimport \"net/http\"\n\ntype LoggingMiddleware struct{}\n\nfunc (m LoggingMiddleware) Wrap(next http.Handler) http.Handler {\n\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n\t\tnext.ServeHTTP(w, r)\n\t})\n}\n"),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/config/config.go", moduleDir),
			Content: []byte("package config\n\ntype Config struct {\n\tPort int\n}\n"),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/config/load.go", moduleDir),
			Content: []byte("package config\n\nfunc Load() Config {\n\treturn Config{Port: 8080}\n}\n"),
		})
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/handler_test.go", moduleDir),
			Content: []byte(fmt.Sprintf("package %s\n\nimport \"testing\"\n\nfunc TestHandler(t *testing.T) {\n\tt.Log(\"placeholder test for %s\")\n}\n", pkg, name)),
		})
	}

	for _, svc := range services {
		var name string
		var port int
		var kind string
		if serviceMap, ok := svc.(map[string]any); ok {
			name = fmt.Sprint(serviceMap["name"])
			if rawPort, ok := serviceMap["port"].(int); ok {
				port = rawPort
			}
			kind = fmt.Sprint(serviceMap["kind"])
		}
		if name == "" {
			continue
		}
		serviceDir := fmt.Sprintf("internal/%s", slugify(name))
		pkg := slugify(name)
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("%s/config.yaml", serviceDir),
			Content: []byte(fmt.Sprintf("name: %s\nport: %d\nkind: %s\n", name, port, kind)),
		})
		if kind == "http" || kind == "" {
			artifacts = append(artifacts, Artifact{
				Path:    fmt.Sprintf("%s/server.go", serviceDir),
				Content: []byte(fmt.Sprintf("package %s\n\nimport \"fmt\"\n\nfunc Run(port int) error {\n\tfmt.Printf(\"%%s listening on :%%d\\n\", %q, port)\n\treturn nil\n}\n", pkg, name)),
			})
			artifacts = append(artifacts, Artifact{
				Path:    fmt.Sprintf("%s/server_test.go", serviceDir),
				Content: []byte(fmt.Sprintf("package %s\n\nimport \"testing\"\n\nfunc TestRun(t *testing.T) {\n\tt.Log(\"placeholder test for %s server\")\n}\n", pkg, name)),
			})
		}
	}

	return artifacts, nil
}

func (DefaultEngine) GenerateForLanguage(neir *model.NEIR, lang language.Language) ([]Artifact, error) {
	if neir == nil {
		return nil, fmt.Errorf("neir is nil")
	}
	if !language.IsValid(lang) {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	projectName := ""
	if neir.Project != nil {
		projectName = neir.Project.Name
	}
	slug := slugify(projectName)

	var artifacts []Artifact

	exts := language.Extensions(lang)
	ext := ".go"
	if len(exts) > 0 {
		ext = exts[0]
	}
	buildFile := language.BuildFile(lang)

	if buildFile != "" {
		artifacts = append(artifacts, Artifact{
			Path:    buildFile,
			Content: []byte(generateBuildFile(lang, slug)),
		})
	}

	artifacts = append(artifacts, Artifact{
		Path:    fmt.Sprintf("src/main%s", ext),
		Content: []byte(generateMainFile(lang, projectName)),
	})

	artifacts = append(artifacts, Artifact{
		Path:    "Dockerfile",
		Content: []byte(generateDockerfile(lang)),
	})

	for _, m := range neir.Modules {
		artifacts = append(artifacts, Artifact{
			Path:    fmt.Sprintf("src/%s/main%s", slugify(m.Name), ext),
			Content: []byte(generateModuleFile(lang, m.Name, projectName)),
		})
	}

	return artifacts, nil
}

func generateBuildFile(lang language.Language, slug string) string {
	switch lang {
	case language.LanguageGo:
		return fmt.Sprintf("module github.com/example/%s\n\ngo 1.22\n", slug)
	case language.LanguageTypeScript:
		return fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "main": "src/main.ts",
  "scripts": {
    "build": "tsc",
    "start": "node dist/main.js"
  }
}
`, slug)
	case language.LanguagePython:
		return fmt.Sprintf(`[project]
name = "%s"
version = "1.0.0"
requires-python = ">=3.10"
`, slug)
	case language.LanguageJava:
		return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<project>
  <groupId>com.example</groupId>
  <artifactId>%s</artifactId>
  <version>1.0.0</version>
</project>
`, slug)
	case language.LanguageRust:
		return fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"
`, slug)
	default:
		return ""
	}
}

func generateMainFile(lang language.Language, projectName string) string {
	switch lang {
	case language.LanguageGo:
		return fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from %s\")\n}\n", projectName)
	case language.LanguageTypeScript:
		return fmt.Sprintf("console.log('hello from %s');\n", projectName)
	case language.LanguagePython:
		return fmt.Sprintf("print('hello from %s')\n", projectName)
	case language.LanguageJava:
		return fmt.Sprintf("public class App {\n    public static void main(String[] args) {\n        System.out.println(\"hello from %s\");\n    }\n}\n", projectName)
	case language.LanguageRust:
		return fmt.Sprintf("fn main() {\n    println!(\"hello from %s\");\n}\n", projectName)
	default:
		return ""
	}
}

func generateDockerfile(lang language.Language) string {
	switch lang {
	case language.LanguageGo:
		return "FROM golang:1.22-alpine\nWORKDIR /app\nCOPY . .\nRUN go build ./src/main.go\nCMD [\"./main\"]\n"
	case language.LanguageTypeScript:
		return "FROM node:20-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nRUN npm run build\nCMD [\"node\", \"dist/main.js\"]\n"
	case language.LanguagePython:
		return "FROM python:3.12-slim\nWORKDIR /app\nCOPY . .\nCMD [\"python\", \"src/main.py\"]\n"
	case language.LanguageJava:
		return "FROM eclipse-temurin:21-jdk\nWORKDIR /app\nCOPY . .\nRUN javac src/main.java\nCMD [\"java\", \"src/App\"]\n"
	case language.LanguageRust:
		return "FROM rust:1.75-slim\nWORKDIR /app\nCOPY . .\nRUN cargo build --release\nCMD [\"./target/release/app\"]\n"
	default:
		return ""
	}
}

func generateModuleFile(lang language.Language, moduleName, projectName string) string {
	pkg := slugify(moduleName)
	switch lang {
	case language.LanguageGo:
		return fmt.Sprintf("package %s\n\n// %s module.\n", pkg, moduleName)
	case language.LanguageTypeScript:
		return fmt.Sprintf("// %s module\nexport {};\n", moduleName)
	case language.LanguagePython:
		return fmt.Sprintf("# %s module\n", moduleName)
	case language.LanguageJava:
		return fmt.Sprintf("public class %s {\n    // %s module\n}\n", strings.Title(moduleName), moduleName)
	case language.LanguageRust:
		return fmt.Sprintf("// %s module\n", moduleName)
	default:
		return ""
	}
}

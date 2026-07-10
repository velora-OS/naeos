package adapters

import (
	"fmt"
	"strings"

	"github.com/NAEOS-foundation/naeos/internal/generation/engine"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
)

type GoAdapter struct{}

func init() {
	Register(GoAdapter{})
}

func (GoAdapter) Language() language.Language {
	return language.LanguageGo
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func pkgName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "/", "")
	return s
}

func (GoAdapter) GenerateProject(projectName string) []engine.Artifact {
	slug := slugify(projectName)
	pkg := pkgName(projectName)

	return []engine.Artifact{
		{Path: "README.md", Content: []byte(fmt.Sprintf("# %s\n\nGenerated from NAEOS pipeline (Go).\n\n## Quick Start\n\n```bash\ngo run ./cmd/app\n```\n\n## Test\n\n```bash\ngo test ./...\n```\n", projectName))},
		{Path: "go.mod", Content: []byte(fmt.Sprintf("module github.com/example/%s\n\ngo 1.22\n", slug))},
		{Path: "cmd/app/main.go", Content: []byte(fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello from %s\")\n}\n", projectName))},
		{Path: fmt.Sprintf("%s/package.go", pkg), Content: []byte(fmt.Sprintf("package %s\n\n// %s module.\n", pkg, projectName))},
	}
}

func (GoAdapter) GenerateModule(moduleName, modulePath, projectName string) []engine.Artifact {
	dir := slugify(modulePath)
	if dir == "" {
		dir = slugify(moduleName)
	}
	pkg := pkgName(moduleName)

	return []engine.Artifact{
		{Path: fmt.Sprintf("%s/handler.go", dir), Content: []byte(fmt.Sprintf("package %s\n\ntype Handler struct {\n\tservice Service\n}\n\nfunc NewHandler(service Service) *Handler {\n\treturn &Handler{service: service}\n}\n", pkg))},
		{Path: fmt.Sprintf("%s/repository.go", dir), Content: []byte(fmt.Sprintf("package %s\n\ntype Repository interface {\n\tList() []string\n}\n", pkg))},
		{Path: fmt.Sprintf("%s/service.go", dir), Content: []byte(fmt.Sprintf("package %s\n\ntype Service interface {\n\tHandle() string\n}\n", pkg))},
		{Path: fmt.Sprintf("%s/domain/model.go", dir), Content: []byte("package domain\n\ntype Model struct {\n\tName string\n}\n")},
		{Path: fmt.Sprintf("%s/http/handler.go", dir), Content: []byte(fmt.Sprintf("package http\n\nimport \"net/http\"\n\ntype Handler struct{}\n\nfunc (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n\tw.Write([]byte(\"handler for %s\"))\n}\n", moduleName))},
		{Path: fmt.Sprintf("%s/http/router.go", dir), Content: []byte("package http\n\ntype Router struct{}\n")},
		{Path: fmt.Sprintf("%s/middleware/logging.go", dir), Content: []byte("package middleware\n\nimport \"net/http\"\n\ntype LoggingMiddleware struct{}\n\nfunc (m LoggingMiddleware) Wrap(next http.Handler) http.Handler {\n\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n\t\tnext.ServeHTTP(w, r)\n\t})\n}\n")},
		{Path: fmt.Sprintf("%s/config/config.go", dir), Content: []byte("package config\n\ntype Config struct {\n\tPort int\n}\n")},
		{Path: fmt.Sprintf("%s/handler_test.go", dir), Content: []byte(fmt.Sprintf("package %s\n\nimport \"testing\"\n\nfunc TestHandler(t *testing.T) {\n\tt.Log(\"test for %s\")\n}\n", pkg, moduleName))},
	}
}

func (GoAdapter) GenerateService(serviceName, serviceKind string, servicePort int, projectName string) []engine.Artifact {
	dir := fmt.Sprintf("internal/%s", slugify(serviceName))
	pkg := pkgName(serviceName)

	var artifacts []engine.Artifact
	artifacts = append(artifacts, engine.Artifact{
		Path:    fmt.Sprintf("%s/config.yaml", dir),
		Content: []byte(fmt.Sprintf("name: %s\nport: %d\nkind: %s\n", serviceName, servicePort, serviceKind)),
	})

	if serviceKind == "http" || serviceKind == "" {
		artifacts = append(artifacts, engine.Artifact{
			Path:    fmt.Sprintf("%s/server.go", dir),
			Content: []byte(fmt.Sprintf("package %s\n\nimport \"fmt\"\n\nfunc Run(port int) error {\n\tfmt.Printf(\"%%s listening on :%%d\\n\", %q, port)\n\treturn nil\n}\n", pkg, serviceName)),
		})
		artifacts = append(artifacts, engine.Artifact{
			Path:    fmt.Sprintf("%s/server_test.go", dir),
			Content: []byte(fmt.Sprintf("package %s\n\nimport \"testing\"\n\nfunc TestRun(t *testing.T) {\n\tt.Log(\"test for %s server\")\n}\n", pkg, serviceName)),
		})
	}

	return artifacts
}

func (GoAdapter) GenerateDockerfile(projectName string) []engine.Artifact {
	return []engine.Artifact{{
		Path:    "Dockerfile",
		Content: []byte("FROM golang:1.22-alpine AS build\nWORKDIR /app\nCOPY . .\nRUN go build ./cmd/app\n\nFROM alpine:3.19\nCOPY --from=build /app/app /app/app\nCMD [\"/app/app\"]\n"),
	}}
}

func (GoAdapter) GenerateCI(projectName string) []engine.Artifact {
	return []engine.Artifact{{
		Path: ".github/workflows/ci.yml",
		Content: []byte("name: ci\n\non: [push, pull_request]\n\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-go@v5\n        with:\n          go-version: '1.22'\n      - run: go test ./...\n"),
	}}
}

func (GoAdapter) GenerateDockerCompose(projectName string) []engine.Artifact {
	return []engine.Artifact{{
		Path:    "docker-compose.yml",
		Content: []byte("version: '3.8'\nservices:\n  app:\n    build: .\n    ports:\n      - '8080:8080'\n"),
	}}
}

func (GoAdapter) GenerateArchitectureDoc(projectName, pattern string) []engine.Artifact {
	return []engine.Artifact{{
		Path:    "docs/architecture.md",
		Content: []byte(fmt.Sprintf("# Architecture\n\nPattern: %s\n\nProject: %s\n", pattern, projectName)),
	}}
}

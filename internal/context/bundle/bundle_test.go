package contextbundle

import (
	"fmt"
	"strings"
	"testing"

	"github.com/NAEOS-foundation/naeos/internal/neir/model"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/generation"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/language"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/module"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/project"
	"github.com/NAEOS-foundation/naeos/internal/neir/model/service"
	"github.com/NAEOS-foundation/naeos/internal/specification/parser"
)

func TestGenerateFromNEIR(t *testing.T) {
	neir := &model.NEIR{
		Project: &project.Project{Name: "test-project"},
		Modules: []module.Module{
			{Name: "auth", Path: "./auth", Dependencies: []string{"core"}},
			{Name: "api", Path: "./api", Dependencies: []string{"auth"}},
		},
		Services: []service.Service{
			{Name: "gateway", Kind: service.KindHTTP, Port: 8080},
		},
		Generation: &generation.GenerationConfig{
			Languages: []language.Language{"go", "typescript"},
		},
	}

	gen := NewGenerator(nil)
	bundle := gen.GenerateFromNEIR(neir)

	if bundle.Project != "test-project" {
		t.Errorf("project = %q, want test-project", bundle.Project)
	}
	if len(bundle.Modules) != 2 {
		t.Errorf("modules = %d, want 2", len(bundle.Modules))
	}
	if bundle.Modules[0].Name != "auth" {
		t.Errorf("modules[0].name = %q, want auth", bundle.Modules[0].Name)
	}
	if bundle.Modules[0].Dependencies[0] != "core" {
		t.Errorf("modules[0].deps[0] = %q, want core", bundle.Modules[0].Dependencies[0])
	}
	if len(bundle.Services) != 1 {
		t.Errorf("services = %d, want 1", len(bundle.Services))
	}
	if bundle.Services[0].Port != 8080 {
		t.Errorf("services[0].port = %d, want 8080", bundle.Services[0].Port)
	}
	if len(bundle.Languages) != 2 {
		t.Errorf("languages = %d, want 2", len(bundle.Languages))
	}
	if bundle.Summary == "" {
		t.Error("summary should not be empty")
	}
	if bundle.Metadata["module_count"] != "2" {
		t.Errorf("metadata module_count = %q, want 2", bundle.Metadata["module_count"])
	}
}

func TestGenerateFromSpec(t *testing.T) {
	doc := &parser.SpecDocument{
		Project: "my-app",
		Modules: []parser.Module{
			{Name: "web", Path: "./web", Description: "web frontend"},
		},
		Services: []parser.Service{
			{Name: "api-server", Kind: "rest", Port: 3000, Endpoints: []parser.Endpoint{
				{Method: "GET", Path: "/users", Action: "listUsers"},
			}},
		},
		Generation: &parser.Generation{Languages: []string{"go"}},
	}

	gen := NewGenerator(nil)
	bundle := gen.GenerateFromSpec(doc)

	if bundle.Project != "my-app" {
		t.Errorf("project = %q, want my-app", bundle.Project)
	}
	if bundle.Modules[0].Description != "web frontend" {
		t.Errorf("description = %q, want web frontend", bundle.Modules[0].Description)
	}
	if bundle.Services[0].Endpoints[0].Action != "listUsers" {
		t.Errorf("endpoint action = %q, want listUsers", bundle.Services[0].Endpoints[0].Action)
	}
}

func TestToMarkdown(t *testing.T) {
	bundle := &Bundle{
		Project:   "test",
		Modules:   []ModuleContext{{Name: "auth", Path: "./auth", Dependencies: []string{"core"}}},
		Services:  []ServiceContext{{Name: "api", Kind: "rest", Port: 8080}},
		Languages: []string{"go"},
		Summary:   "Project: test",
	}

	md := bundle.ToMarkdown()
	if !strings.Contains(md, "# test") {
		t.Error("markdown should contain project title")
	}
	if !strings.Contains(md, "## Modules") {
		t.Error("markdown should contain modules section")
	}
	if !strings.Contains(md, "auth") {
		t.Error("markdown should contain module name")
	}
	if !strings.Contains(md, "core") {
		t.Error("markdown should contain dependency")
	}
}

func TestToPlainText(t *testing.T) {
	bundle := &Bundle{
		Project:  "test",
		Modules:  []ModuleContext{{Name: "auth", Path: "./auth"}},
		Services: []ServiceContext{{Name: "api", Kind: "rest", Port: 8080}},
	}

	plain := bundle.ToPlainText()
	if !strings.Contains(plain, "Project: test") {
		t.Error("plain text should contain project")
	}
	if !strings.Contains(plain, "Module: auth") {
		t.Error("plain text should contain module")
	}
}

func TestSupportedTargets(t *testing.T) {
	bundle := &Bundle{
		Modules: []ModuleContext{{Name: "x"}},
	}

	targets := bundle.SupportedTargets()
	found := false
	for _, tgt := range targets {
		if tgt == "markdown" {
			found = true
		}
	}
	if !found {
		t.Error("markdown should be a supported target")
	}
}

func TestBundleMetadata(t *testing.T) {
	bundle := &Bundle{
		Project:  "meta-test",
		Metadata: make(map[string]string),
		Modules:  []ModuleContext{{Name: "a"}, {Name: "b"}},
		Services: []ServiceContext{{Name: "s1"}},
	}
	bundle.Metadata["module_count"] = fmt.Sprintf("%d", len(bundle.Modules))
	bundle.Metadata["service_count"] = fmt.Sprintf("%d", len(bundle.Services))

	if bundle.Metadata["module_count"] != "2" {
		t.Errorf("module_count = %q, want 2", bundle.Metadata["module_count"])
	}
	if bundle.Metadata["service_count"] != "1" {
		t.Errorf("service_count = %q, want 1", bundle.Metadata["service_count"])
	}
}

func TestGenerateFromNEIREmpty(t *testing.T) {
	neir := &model.NEIR{}

	gen := NewGenerator(nil)
	bundle := gen.GenerateFromNEIR(neir)

	if bundle.Project != "" {
		t.Errorf("project should be empty, got %q", bundle.Project)
	}
	if len(bundle.Modules) != 0 {
		t.Errorf("modules should be empty, got %d", len(bundle.Modules))
	}
	_ = bundle.Summary
}

func TestEstimateTokens(t *testing.T) {
	bundle := &Bundle{
		Project:  "test",
		Modules:  []ModuleContext{{Name: "auth"}, {Name: "api"}},
		Services: []ServiceContext{{Name: "srv"}},
		Languages: []string{"go"},
	}
	tokens := bundle.EstimateTokens()
	if tokens <= 0 {
		t.Errorf("tokens = %d, want > 0", tokens)
	}
}

func TestToJSON(t *testing.T) {
	bundle := &Bundle{
		Project:  "json-test",
		Modules:  []ModuleContext{{Name: "m1"}},
		Metadata: map[string]string{"k": "v"},
	}
	out := bundle.ToJSON()
	if !strings.Contains(out, "json-test") {
		t.Error("JSON should contain project name")
	}
	if !strings.Contains(out, "m1") {
		t.Error("JSON should contain module name")
	}
}

func TestToJSONEmpty(t *testing.T) {
	bundle := &Bundle{}
	out := bundle.ToJSON()
	if !strings.Contains(out, "project") {
		t.Error("empty bundle JSON should contain field names")
	}
}

func TestFilterByModule(t *testing.T) {
	bundle := &Bundle{
		Project: "proj",
		Modules: []ModuleContext{
			{Name: "auth"},
			{Name: "api"},
			{Name: "core"},
		},
		DependencyGraph: []DependencyEdge{
			{From: "auth", To: "core"},
			{From: "api", To: "auth"},
		},
		Services: []ServiceContext{{Name: "s1"}},
	}
	filtered := bundle.FilterByModule([]string{"auth", "api"})
	if len(filtered.Modules) != 2 {
		t.Errorf("modules = %d, want 2", len(filtered.Modules))
	}
	if len(filtered.DependencyGraph) != 2 {
		t.Errorf("dep graph = %d, want 2", len(filtered.DependencyGraph))
	}
	if len(filtered.Services) != 1 {
		t.Error("services should be preserved")
	}
}

func TestFilterByModuleEmpty(t *testing.T) {
	bundle := &Bundle{
		Modules: []ModuleContext{{Name: "a"}, {Name: "b"}},
	}
	filtered := bundle.FilterByModule([]string{"nonexistent"})
	if len(filtered.Modules) != 0 {
		t.Errorf("modules = %d, want 0", len(filtered.Modules))
	}
}

func TestFilterByService(t *testing.T) {
	bundle := &Bundle{
		Modules:  []ModuleContext{{Name: "m1"}},
		Services: []ServiceContext{
			{Name: "http-svc", Kind: "rest"},
			{Name: "grpc-svc", Kind: "grpc"},
			{Name: "ws-svc", Kind: "websocket"},
		},
	}
	filtered := bundle.FilterByService([]string{"rest", "grpc"})
	if len(filtered.Services) != 2 {
		t.Errorf("services = %d, want 2", len(filtered.Services))
	}
	if len(filtered.Modules) != 1 {
		t.Error("modules should be preserved")
	}
}

func TestFilterByServiceCaseInsensitive(t *testing.T) {
	bundle := &Bundle{
		Services: []ServiceContext{
			{Name: "s1", Kind: "REST"},
		},
	}
	filtered := bundle.FilterByService([]string{"rest"})
	if len(filtered.Services) != 1 {
		t.Errorf("services = %d, want 1 (case insensitive)", len(filtered.Services))
	}
}

func TestMergeWithOther(t *testing.T) {
	b1 := &Bundle{
		Project:  "p1",
		Modules:  []ModuleContext{{Name: "m1"}},
		Services: []ServiceContext{{Name: "s1"}},
		Languages: []string{"go"},
		Metadata:  map[string]string{"a": "1"},
	}
	b2 := &Bundle{
		Project:  "p2",
		Modules:  []ModuleContext{{Name: "m2"}},
		Services: []ServiceContext{{Name: "s2"}},
		Languages: []string{"typescript"},
		Metadata:  map[string]string{"b": "2"},
		Security:  &SecurityContext{AuthMethod: "jwt"},
	}
	merged := b1.Merge(b2)
	if merged.Project != "p2" {
		t.Errorf("project = %q, want p2", merged.Project)
	}
	if len(merged.Modules) != 2 {
		t.Errorf("modules = %d, want 2", len(merged.Modules))
	}
	if len(merged.Services) != 2 {
		t.Errorf("services = %d, want 2", len(merged.Services))
	}
	if len(merged.Languages) != 2 {
		t.Errorf("languages = %d, want 2", len(merged.Languages))
	}
	if merged.Security == nil || merged.Security.AuthMethod != "jwt" {
		t.Error("security should come from other")
	}
	if merged.Metadata["a"] != "1" || merged.Metadata["b"] != "2" {
		t.Error("metadata should be merged")
	}
	if merged.Metadata["token_estimate"] == "" {
		t.Error("token_estimate should be set")
	}
	if merged.Summary == "" {
		t.Error("summary should be set")
	}
}

func TestMergeNilOther(t *testing.T) {
	b := &Bundle{
		Project:  "p1",
		Modules:  []ModuleContext{{Name: "m1"}},
		Languages: []string{"go"},
		Metadata: map[string]string{"k": "v"},
	}
	merged := b.Merge(nil)
	if merged.Project != "p1" {
		t.Errorf("project = %q, want p1", merged.Project)
	}
	if len(merged.Modules) != 1 {
		t.Errorf("modules = %d, want 1", len(merged.Modules))
	}
}

func TestMergeDedupLanguages(t *testing.T) {
	b1 := &Bundle{Languages: []string{"go", "typescript"}}
	b2 := &Bundle{Languages: []string{"go", "python"}}
	merged := b1.Merge(b2)
	if len(merged.Languages) != 3 {
		t.Errorf("languages = %d, want 3 (deduped)", len(merged.Languages))
	}
}

func TestMergeModulesFromBoth(t *testing.T) {
	b1 := &Bundle{
		Modules: []ModuleContext{{Name: "m1"}},
	}
	b2 := &Bundle{
		Modules: []ModuleContext{{Name: "m2"}},
	}
	merged := b1.Merge(b2)
	if len(merged.Modules) != 2 {
		t.Errorf("modules = %d, want 2", len(merged.Modules))
	}
}

func TestMergeNoModules(t *testing.T) {
	b1 := &Bundle{Project: "p1", Metadata: map[string]string{}}
	b2 := &Bundle{Project: "p2", Metadata: map[string]string{}}
	merged := b1.Merge(b2)
	if len(merged.Modules) != 0 {
		t.Errorf("modules = %d, want 0", len(merged.Modules))
	}
}

func TestSupportedTargetsWithNEIR(t *testing.T) {
	bundle := &Bundle{NEIR: "some json"}
	targets := bundle.SupportedTargets()
	found := false
	for _, tgt := range targets {
		if tgt == "neir" {
			found = true
		}
	}
	if !found {
		t.Error("neir should be a supported target when NEIR is set")
	}
}

func TestBuildTargets(t *testing.T) {
	bundle := &Bundle{}
	targets := bundle.buildTargets([]string{"go", "typescript"})
	hasLangGo := false
	hasLangTS := false
	for _, tgt := range targets {
		if tgt == "lang-go" {
			hasLangGo = true
		}
		if tgt == "lang-typescript" {
			hasLangTS = true
		}
	}
	if !hasLangGo {
		t.Error("should have lang-go target")
	}
	if !hasLangTS {
		t.Error("should have lang-typescript target")
	}
}

func TestBuildSummary(t *testing.T) {
	bundle := &Bundle{
		Project:  "proj",
		Modules:  []ModuleContext{{Name: "a"}, {Name: "b"}},
		Services: []ServiceContext{{Name: "s1"}},
		Languages: []string{"go"},
	}
	summary := bundle.buildSummary()
	if !strings.Contains(summary, "Project: proj") {
		t.Error("summary should contain project")
	}
	if !strings.Contains(summary, "Modules: a, b") {
		t.Error("summary should contain modules")
	}
	if !strings.Contains(summary, "Services: 1") {
		t.Error("summary should contain service count")
	}
	if !strings.Contains(summary, "Languages: go") {
		t.Error("summary should contain languages")
	}
}

func TestBuildSummaryEmpty(t *testing.T) {
	bundle := &Bundle{}
	summary := bundle.buildSummary()
	if summary != "" {
		t.Errorf("empty bundle summary = %q, want empty", summary)
	}
}

func TestFilterByModulePreservesCloudAndSecurity(t *testing.T) {
	bundle := &Bundle{
		Modules: []ModuleContext{{Name: "m1"}},
		Cloud:   []CloudResource{{Provider: "aws", Type: "s3", Name: "bucket"}},
		Security: &SecurityContext{AuthMethod: "oauth2"},
	}
	filtered := bundle.FilterByModule([]string{"m1"})
	if filtered.Cloud == nil || len(filtered.Cloud) != 1 {
		t.Error("cloud resources should be preserved")
	}
	if filtered.Security == nil || filtered.Security.AuthMethod != "oauth2" {
		t.Error("security should be preserved")
	}
}

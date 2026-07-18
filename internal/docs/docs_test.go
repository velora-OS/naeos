package docs

import (
	"strings"
	"testing"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func lines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

// --- DocGenerator tests ---

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestNewGeneratorWithArtifacts(t *testing.T) {
	artifacts := []ArtifactRef{
		{Path: "main.go", Size: 100, Type: "go"},
	}
	gen := NewGenerator("proj", artifacts)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestGenerateAPIDocs(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	endpoints := []Endpoint{
		{Method: "GET", Path: "/health", Description: "Health check"},
		{Method: "POST", Path: "/api/v1/users", Description: "Create user"},
	}
	content := gen.GenerateAPIDocs(endpoints)
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "/health") {
		t.Error("expected /health endpoint in docs")
	}
}

func TestGenerateAPIDocsEmpty(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	content := gen.GenerateAPIDocs([]Endpoint{})
	if !contains(content, "API Documentation") {
		t.Error("expected header even with no endpoints")
	}
	if !contains(content, "| Method |") {
		t.Error("expected table header")
	}
}

func TestGenerateAPIDocsProjectName(t *testing.T) {
	gen := NewGenerator("my-api", nil)
	content := gen.GenerateAPIDocs([]Endpoint{
		{Method: "GET", Path: "/", Description: "root"},
	})
	if !contains(content, "# my-api API Documentation") {
		t.Error("expected project name in header")
	}
}

func TestGenerateArchitectureDiagram(t *testing.T) {
	gen := NewGenerator("test-project", nil)
	content := gen.GenerateArchitectureDiagram(
		[]string{"api", "worker"},
		[]string{"core", "auth"},
	)
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "mermaid") {
		t.Error("expected mermaid diagram")
	}
}

func TestGenerateArchitectureDiagramNoServices(t *testing.T) {
	gen := NewGenerator("proj", nil)
	content := gen.GenerateArchitectureDiagram([]string{}, []string{"core"})
	if !contains(content, "core") {
		t.Error("expected core module in diagram")
	}
}

func TestGenerateArchitectureDiagramNoModules(t *testing.T) {
	gen := NewGenerator("proj", nil)
	content := gen.GenerateArchitectureDiagram([]string{"api"}, []string{})
	if !contains(content, "api") {
		t.Error("expected api service in diagram")
	}
}

func TestGenerateProjectDocs(t *testing.T) {
	gen := NewGenerator("test-project", []ArtifactRef{
		{Path: "main.go", Size: 100, Type: "go"},
	})
	content := gen.GenerateProjectDocs()
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !contains(content, "test-project") {
		t.Error("expected project name in docs")
	}
}

func TestGenerateProjectDocsMultipleArtifacts(t *testing.T) {
	gen := NewGenerator("proj", []ArtifactRef{
		{Path: "a.go", Size: 50, Type: "go"},
		{Path: "b.go", Size: 200, Type: "go"},
	})
	content := gen.GenerateProjectDocs()
	if !contains(content, "a.go") || !contains(content, "b.go") {
		t.Error("expected both artifacts in docs")
	}
}

func TestGenerateProjectDocsEmpty(t *testing.T) {
	gen := NewGenerator("proj", nil)
	content := gen.GenerateProjectDocs()
	if !contains(content, "Quick Start") {
		t.Error("expected Quick Start section")
	}
}

// --- ChangelogGenerator tests ---

func TestChangelogGeneratorGenerate(t *testing.T) {
	cg := NewChangelogGenerator("my-app", []ChangeEntry{
		{
			Version: "1.0.0",
			Date:    "2024-01-15",
			Changes: []ChangeItem{
				{Category: "Added", Text: "Initial release"},
				{Category: "Added", Text: "Documentation"},
				{Category: "Fixed", Text: "Bug in login"},
			},
		},
	})
	content := cg.Generate()
	if !contains(content, "## 1.0.0 (2024-01-15)") {
		t.Error("expected version header")
	}
	if !contains(content, "### Added") {
		t.Error("expected Added category")
	}
	if !contains(content, "- Initial release") {
		t.Error("expected change item")
	}
	if !contains(content, "### Fixed") {
		t.Error("expected Fixed category")
	}
}

func TestChangelogGeneratorLatestVersion(t *testing.T) {
	cg := NewChangelogGenerator("app", []ChangeEntry{
		{Version: "2.0.0", Date: "2024-06-01"},
		{Version: "1.0.0", Date: "2024-01-01"},
	})
	if cg.LatestVersion() != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", cg.LatestVersion())
	}
}

func TestChangelogGeneratorLatestVersionEmpty(t *testing.T) {
	cg := NewChangelogGenerator("app", nil)
	if cg.LatestVersion() != "" {
		t.Error("expected empty for no entries")
	}
}

func TestChangelogGeneratorEntriesForVersion(t *testing.T) {
	cg := NewChangelogGenerator("app", []ChangeEntry{
		{Version: "1.0.0", Changes: []ChangeItem{{Category: "Added", Text: "feat a"}}},
		{Version: "1.1.0", Changes: []ChangeItem{{Category: "Added", Text: "feat b"}}},
	})
	entry := cg.EntriesForVersion("1.1.0")
	if entry == nil {
		t.Fatal("expected entry for 1.1.0")
	}
	if entry.Version != "1.1.0" {
		t.Errorf("expected version 1.1.0, got %s", entry.Version)
	}
}

func TestChangelogGeneratorEntriesForVersionNotFound(t *testing.T) {
	cg := NewChangelogGenerator("app", []ChangeEntry{
		{Version: "1.0.0"},
	})
	if cg.EntriesForVersion("9.0.0") != nil {
		t.Error("expected nil for missing version")
	}
}

func TestChangelogGeneratorMultipleVersions(t *testing.T) {
	cg := NewChangelogGenerator("app", []ChangeEntry{
		{Version: "1.0.0", Date: "2024-01-01", Changes: []ChangeItem{
			{Category: "Added", Text: "feat a"},
			{Category: "Deprecated", Text: "old API"},
		}},
		{Version: "0.9.0", Date: "2023-12-01", Changes: []ChangeItem{
			{Category: "Added", Text: "initial"},
		}},
	})
	content := cg.Generate()
	if !contains(content, "1.0.0") || !contains(content, "0.9.0") {
		t.Error("expected both versions")
	}
	if !contains(content, "### Deprecated") {
		t.Error("expected Deprecated category")
	}
}

// --- ContributorGuide tests ---

func TestContributorGuideGenerateDefault(t *testing.T) {
	cg := NewContributorGuide("my-app", "")
	content := cg.Generate()
	if !contains(content, "# Contributing to my-app") {
		t.Error("expected title")
	}
	if !contains(content, "Fork the repository") {
		t.Error("expected default getting started section")
	}
	if !contains(content, "Code of Conduct") {
		t.Error("expected default code of conduct")
	}
}

func TestContributorGuideWithRepoURL(t *testing.T) {
	cg := NewContributorGuide("app", "https://github.com/org/app")
	content := cg.Generate()
	if !contains(content, "https://github.com/org/app") {
		t.Error("expected repo URL in output")
	}
}

func TestContributorGuideCustomSections(t *testing.T) {
	cg := NewContributorGuide("app", "")
	cg.AddSection("Prerequisites", "Install Go 1.21+\n")
	cg.AddSection("Running Tests", "Run `go test ./...`\n")
	content := cg.Generate()
	if !contains(content, "## Prerequisites") {
		t.Error("expected Prerequisites section")
	}
	if !contains(content, "Install Go 1.21+") {
		t.Error("expected prerequisites content")
	}
	if !contains(content, "## Running Tests") {
		t.Error("expected Running Tests section")
	}
	if contains(content, "Fork the repository") {
		t.Error("default sections should not appear when custom sections exist")
	}
}

func TestContributorGuideSectionCount(t *testing.T) {
	cg := NewContributorGuide("app", "")
	if cg.SectionCount() != 0 {
		t.Errorf("expected 0, got %d", cg.SectionCount())
	}
	cg.AddSection("a", "b")
	cg.AddSection("c", "d")
	if cg.SectionCount() != 2 {
		t.Errorf("expected 2, got %d", cg.SectionCount())
	}
}

// --- ConfigDoc tests ---

func TestConfigDocGenerate(t *testing.T) {
	cd := NewConfigDoc("my-app", map[string]ConfigField{
		"server.port": {
			Name: "server.port", Type: "int", Default: "8080",
			Description: "Server port", Required: false,
		},
		"database.host": {
			Name: "database.host", Type: "string", Default: "localhost",
			Description: "Database host", Required: true,
		},
	})
	content := cd.Generate()
	if !contains(content, "# my-app Configuration") {
		t.Error("expected config header")
	}
	if !contains(content, "server.port") {
		t.Error("expected server.port field")
	}
	if !contains(content, "database.host") {
		t.Error("expected database.host field")
	}
	if !contains(content, "Environment Variables") {
		t.Error("expected environment variables section")
	}
}

func TestConfigDocFieldsSorted(t *testing.T) {
	cd := NewConfigDoc("app", map[string]ConfigField{
		"zebra": {Name: "zebra", Type: "string"},
		"alpha": {Name: "alpha", Type: "string"},
	})
	content := cd.Generate()
	alphaIdx := strings.Index(content, "alpha")
	zebraIdx := strings.Index(content, "zebra")
	if alphaIdx > zebraIdx {
		t.Error("expected alpha before zebra (sorted)")
	}
}

func TestConfigDocRequiredFields(t *testing.T) {
	cd := NewConfigDoc("app", map[string]ConfigField{
		"a": {Name: "a", Required: false},
		"b": {Name: "b", Required: true},
		"c": {Name: "c", Required: true},
	})
	required := cd.RequiredFields()
	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}
}

func TestConfigDocFieldCount(t *testing.T) {
	cd := NewConfigDoc("app", map[string]ConfigField{
		"x": {Name: "x"},
		"y": {Name: "y"},
		"z": {Name: "z"},
	})
	if cd.FieldCount() != 3 {
		t.Errorf("expected 3, got %d", cd.FieldCount())
	}
}

func TestConfigDocEnvVarMapping(t *testing.T) {
	cd := NewConfigDoc("app", map[string]ConfigField{
		"server.port": {Name: "server.port", Type: "int"},
	})
	content := cd.Generate()
	if !contains(content, "SERVER_PORT") {
		t.Error("expected env var SERVER_PORT mapped from server.port")
	}
}

func TestConfigDocEmpty(t *testing.T) {
	cd := NewConfigDoc("app", map[string]ConfigField{})
	content := cd.Generate()
	if !contains(content, "Configuration") {
		t.Error("expected header even with no fields")
	}
	if cd.FieldCount() != 0 {
		t.Error("expected 0 fields")
	}
}

// --- MarkdownSection tests ---

func TestMarkdownSectionRender(t *testing.T) {
	ms := NewMarkdownSection("My Title", 2)
	ms.WriteParagraph("Some text here.")
	ms.WriteList([]string{"item 1", "item 2"})
	content := ms.Render()
	if !contains(content, "## My Title") {
		t.Error("expected h2 title")
	}
	if !contains(content, "Some text here.") {
		t.Error("expected paragraph content")
	}
	if !contains(content, "- item 1") {
		t.Error("expected list items")
	}
}

func TestMarkdownSectionWriteCodeBlock(t *testing.T) {
	ms := NewMarkdownSection("Code", 3)
	ms.WriteCodeBlock("go", "fmt.Println(\"hello\")")
	content := ms.Render()
	if !contains(content, "```go") {
		t.Error("expected go code block")
	}
	if !contains(content, "fmt.Println") {
		t.Error("expected code content")
	}
}

func TestMarkdownSectionWriteTable(t *testing.T) {
	ms := NewMarkdownSection("Data", 2)
	ms.WriteTable(
		[]string{"Name", "Value"},
		[][]string{
			{"a", "1"},
			{"b", "2"},
		},
	)
	content := ms.Render()
	if !contains(content, "| Name | Value |") {
		t.Error("expected table header")
	}
	if !contains(content, "| a | 1 |") {
		t.Error("expected table row")
	}
}

func TestMarkdownSectionLevelClamping(t *testing.T) {
	ms1 := NewMarkdownSection("Title", 0)
	if ms1.Render() == "" {
		t.Error("expected non-empty")
	}
	if !contains(ms1.Render(), "# ") {
		t.Error("expected level 1 (clamped from 0)")
	}

	ms2 := NewMarkdownSection("Title", 10)
	if !contains(ms2.Render(), "###### ") {
		t.Error("expected level 6 (clamped from 10)")
	}
}

func TestMarkdownSectionEmptyContent(t *testing.T) {
	ms := NewMarkdownSection("Empty", 1)
	content := ms.Render()
	if !contains(content, "# Empty") {
		t.Error("expected title")
	}
	if !contains(content, "\n\n") {
		t.Error("expected empty content between title and end")
	}
}

// --- MarkdownDocument tests ---

func TestMarkdownDocumentRender(t *testing.T) {
	doc := NewMarkdownDocument("My Doc")
	s1 := NewMarkdownSection("Intro", 1)
	s1.WriteParagraph("Welcome.")
	s2 := NewMarkdownSection("Details", 2)
	s2.WriteList([]string{"point one", "point two"})
	doc.AddSection(s1)
	doc.AddSection(s2)
	content := doc.Render()
	if !contains(content, "# My Doc") {
		t.Error("expected document title")
	}
	if !contains(content, "Welcome.") {
		t.Error("expected intro content")
	}
	if !contains(content, "point one") {
		t.Error("expected details content")
	}
}

func TestMarkdownDocumentSectionCount(t *testing.T) {
	doc := NewMarkdownDocument("doc")
	if doc.SectionCount() != 0 {
		t.Error("expected 0 sections")
	}
	doc.AddSection(NewMarkdownSection("a", 1))
	doc.AddSection(NewMarkdownSection("b", 2))
	if doc.SectionCount() != 2 {
		t.Errorf("expected 2, got %d", doc.SectionCount())
	}
}

func TestMarkdownDocumentEmpty(t *testing.T) {
	doc := NewMarkdownDocument("Empty Doc")
	content := doc.Render()
	if !contains(content, "# Empty Doc") {
		t.Error("expected title")
	}
}

// --- ReadmeGenerator tests ---

func TestReadmeGeneratorBasic(t *testing.T) {
	rg := NewReadmeGenerator("my-app", "A great app")
	content := rg.Generate()
	if !contains(content, "# my-app") {
		t.Error("expected title")
	}
	if !contains(content, "> A great app") {
		t.Error("expected description")
	}
	if !contains(content, "## Features") {
		t.Error("expected features section")
	}
}

func TestReadmeGeneratorWithRepoURL(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc").
		WithRepoURL("https://github.com/org/app")
	content := rg.Generate()
	if !contains(content, "CONTRIBUTING.md") {
		t.Error("expected contributing link")
	}
	if !contains(content, "badge") || !contains(content, "shields.io") {
		t.Error("expected badges")
	}
}

func TestReadmeGeneratorWithLicense(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc").
		WithLicense("MIT")
	content := rg.Generate()
	if !contains(content, "**License:** MIT") {
		t.Error("expected license")
	}
}

func TestReadmeGeneratorWithLicenseAndRepo(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc").
		WithRepoURL("https://github.com/org/app").
		WithLicense("MIT")
	content := rg.Generate()
	if !contains(content, "license-MIT") {
		t.Error("expected license badge")
	}
}

func TestReadmeGeneratorWithInstallAndUsage(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc").
		WithInstall("go install github.com/org/app@latest").
		WithUsage("app serve --port 8080")
	content := rg.Generate()
	if !contains(content, "## Installation") {
		t.Error("expected installation section")
	}
	if !contains(content, "go install") {
		t.Error("expected install command")
	}
	if !contains(content, "## Usage") {
		t.Error("expected usage section")
	}
	if !contains(content, "app serve") {
		t.Error("expected usage command")
	}
}

func TestReadmeGeneratorWithFeatures(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc").
		WithFeatures([]string{"Fast", "Secure", "Scalable"})
	content := rg.Generate()
	if !contains(content, "- Fast") {
		t.Error("expected Fast feature")
	}
	if !contains(content, "- Secure") {
		t.Error("expected Secure feature")
	}
	if !contains(content, "- Scalable") {
		t.Error("expected Scalable feature")
	}
	if contains(content, "TODO") {
		t.Error("should not contain TODO when features provided")
	}
}

func TestReadmeGeneratorNoFeatures(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc")
	content := rg.Generate()
	if !contains(content, "TODO") {
		t.Error("expected TODO placeholder when no features")
	}
}

func TestReadmeGeneratorNoRepoURL(t *testing.T) {
	rg := NewReadmeGenerator("app", "Desc")
	content := rg.Generate()
	if contains(content, "shields.io") {
		t.Error("expected no badges without repo URL")
	}
	if contains(content, "CONTRIBUTING.md") {
		t.Error("expected no contributing link without repo URL")
	}
}

func TestReadmeGeneratorFullBuilder(t *testing.T) {
	rg := NewReadmeGenerator("full-app", "Full featured app").
		WithRepoURL("https://github.com/org/full-app").
		WithLicense("Apache-2.0").
		WithInstall("make install").
		WithUsage("make run").
		WithFeatures([]string{"Feature A", "Feature B"})
	content := rg.Generate()
	if !contains(content, "# full-app") {
		t.Error("expected title")
	}
	if !contains(content, "make install") {
		t.Error("expected install command")
	}
	if !contains(content, "make run") {
		t.Error("expected usage command")
	}
	if !contains(content, "Apache-2.0") {
		t.Error("expected license")
	}
	if !contains(content, "Feature A") {
		t.Error("expected feature A")
	}
}

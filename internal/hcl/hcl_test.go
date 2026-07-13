package hcl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSimple(t *testing.T) {
	input := []byte(`
project "myapp" {
  version     = "1.0.0"
  description = "My application"
}

service "api" {
  image = "myapp-api"
  port  = 8080
  type  = "backend"
}

infra "infra" {
  engine = "docker"
}
`)

	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "myapp" {
		t.Errorf("expected project name 'myapp', got %q", spec.Project.Name)
	}
	if spec.Project.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", spec.Project.Version)
	}
	if spec.Project.Description != "My application" {
		t.Errorf("expected description 'My application', got %q", spec.Project.Description)
	}
	if len(spec.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(spec.Services))
	}
	svc := spec.Services["api"]
	if svc.Image != "myapp-api" {
		t.Errorf("expected image 'myapp-api', got %q", svc.Image)
	}
	if svc.Port != 8080 {
		t.Errorf("expected port 8080, got %d", svc.Port)
	}
	if svc.Type != "backend" {
		t.Errorf("expected type 'backend', got %q", svc.Type)
	}
	if spec.Infra.Engine != "docker" {
		t.Errorf("expected engine 'docker', got %q", spec.Infra.Engine)
	}
}

func TestParseMultipleServices(t *testing.T) {
	input := []byte(`
project "multi" {
  version = "2.0.0"
}

service "api" {
  port = 8080
  type = "backend"
}

service "web" {
  port = 3000
  type = "frontend"
}

service "worker" {
  port = 9090
  type = "job"
}
`)

	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(spec.Services))
	}
	for _, name := range []string{"api", "web", "worker"} {
		if _, ok := spec.Services[name]; !ok {
			t.Errorf("missing service %q", name)
		}
	}
}

func TestParseComments(t *testing.T) {
	input := []byte(`
# This is a comment
// This is also a comment
project "commented" {
  version = "1.0.0"
}
`)

	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "commented" {
		t.Errorf("expected project 'commented', got %q", spec.Project.Name)
	}
}

func TestParseEmpty(t *testing.T) {
	spec, err := Parse([]byte(""), "empty.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "" {
		t.Errorf("expected empty project name, got %q", spec.Project.Name)
	}
	if len(spec.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(spec.Services))
	}
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.hcl")
	content := `
project "filetest" {
  version = "3.0.0"
}

service "backend" {
  port = 5000
  type = "api"
}
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "filetest" {
		t.Errorf("expected 'filetest', got %q", spec.Project.Name)
	}
	if spec.Project.Version != "3.0.0" {
		t.Errorf("expected '3.0.0', got %q", spec.Project.Version)
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/file.hcl")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseInvalid(t *testing.T) {
	input := []byte(`
project "bad" {
  version = "1.0.0"
  unknown_field = "value"
  broken
}
`)

	_, err := Parse(input, "bad.hcl")
	if err != nil {
		t.Logf("got error (expected for malformed HCL): %v", err)
	}
}

func TestParseProjectOnly(t *testing.T) {
	input := []byte(`
project "minimal" {
  version = "1.0.0"
}
`)
	spec, err := Parse(input, "minimal.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "minimal" {
		t.Errorf("expected 'minimal', got %q", spec.Project.Name)
	}
	if len(spec.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(spec.Services))
	}
}

func TestParseInfraOnly(t *testing.T) {
	input := []byte(`
infra "infra" {
  engine = "kubernetes"
}
`)
	spec, err := Parse(input, "infra.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Infra.Engine != "kubernetes" {
		t.Errorf("expected 'kubernetes', got %q", spec.Infra.Engine)
	}
	if spec.Project.Name != "" {
		t.Errorf("expected empty project name, got %q", spec.Project.Name)
	}
}

func TestParseServiceWithoutImage(t *testing.T) {
	input := []byte(`
project "noimg" {
  version = "1.0.0"
}

service "bare" {
  port = 3000
  type = "frontend"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	svc := spec.Services["bare"]
	if svc.Image != "" {
		t.Errorf("expected empty image, got %q", svc.Image)
	}
	if svc.Port != 3000 {
		t.Errorf("expected port 3000, got %d", svc.Port)
	}
}

func TestParseServiceWithoutPort(t *testing.T) {
	input := []byte(`
service "noport" {
  image = "nginx"
  type  = "static"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	svc := spec.Services["noport"]
	if svc.Port != 0 {
		t.Errorf("expected port 0, got %d", svc.Port)
	}
	if svc.Image != "nginx" {
		t.Errorf("expected image 'nginx', got %q", svc.Image)
	}
}

func TestParseServiceWithoutType(t *testing.T) {
	input := []byte(`
service "notype" {
  image = "redis"
  port  = 6379
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	svc := spec.Services["notype"]
	if svc.Type != "" {
		t.Errorf("expected empty type, got %q", svc.Type)
	}
}

func TestParseServiceInvalidPort(t *testing.T) {
	input := []byte(`
service "badport" {
  port = notanumber
  type = "backend"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	svc := spec.Services["badport"]
	if svc.Port != 0 {
		t.Errorf("expected port 0 for invalid value, got %d", svc.Port)
	}
}

func TestParseMultipleServicesSameType(t *testing.T) {
	input := []byte(`
service "svc1" {
  port = 8001
  type = "backend"
}

service "svc2" {
  port = 8002
  type = "backend"
}

service "svc3" {
  port = 8003
  type = "backend"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Services) != 3 {
		t.Fatalf("expected 3 services, got %d", len(spec.Services))
	}
	for _, name := range []string{"svc1", "svc2", "svc3"} {
		if _, ok := spec.Services[name]; !ok {
			t.Errorf("missing service %q", name)
		}
	}
}

func TestParseCommentsInsideBlock(t *testing.T) {
	input := []byte(`
project "commented" {
  # version comment
  version = "1.0.0"
  // description comment
  description = "Has comments"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Version != "1.0.0" {
		t.Errorf("expected '1.0.0', got %q", spec.Project.Version)
	}
	if spec.Project.Description != "Has comments" {
		t.Errorf("expected 'Has comments', got %q", spec.Project.Description)
	}
}

func TestParseBlankLinesOnly(t *testing.T) {
	input := []byte(`



`)
	spec, err := Parse(input, "blank.hcl")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Project.Name != "" {
		t.Errorf("expected empty project name, got %q", spec.Project.Name)
	}
	if len(spec.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(spec.Services))
	}
}

func TestToYAMLMinimal(t *testing.T) {
	spec := &Spec{
		Project: Project{Name: "minimal"},
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "name: minimal") {
		t.Errorf("expected 'name: minimal' in output, got %q", s)
	}
}

func TestToYAMLFull(t *testing.T) {
	spec := &Spec{
		Project: Project{
			Name:        "full",
			Version:     "2.0.0",
			Description: "Full spec",
		},
		Services: map[string]Service{
			"api": {Image: "myimg", Port: 8080, Type: "backend"},
		},
		Infra: Infra{Engine: "docker"},
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "name: full") {
		t.Errorf("expected 'name: full'")
	}
	if !strings.Contains(s, "version: 2.0.0") {
		t.Errorf("expected 'version: 2.0.0'")
	}
	if !strings.Contains(s, "description: Full spec") {
		t.Errorf("expected 'description: Full spec'")
	}
	if !strings.Contains(s, "image: myimg") {
		t.Errorf("expected 'image: myimg'")
	}
	if !strings.Contains(s, "port: 8080") {
		t.Errorf("expected 'port: 8080'")
	}
	if !strings.Contains(s, "type: backend") {
		t.Errorf("expected 'type: backend'")
	}
	if !strings.Contains(s, "engine: docker") {
		t.Errorf("expected 'engine: docker'")
	}
}

func TestToYAMLNoServicesNoInfra(t *testing.T) {
	spec := &Spec{
		Project: Project{Name: "simple", Version: "1.0.0"},
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if strings.Contains(s, "services:") {
		t.Error("should not contain 'services:' when no services")
	}
	if strings.Contains(s, "infra:") {
		t.Error("should not contain 'infra:' when no infra")
	}
}

func TestToYAMLServiceWithoutOptionalFields(t *testing.T) {
	spec := &Spec{
		Project: Project{Name: "bare"},
		Services: map[string]Service{
			"worker": {Type: "job"},
		},
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "name: worker") {
		t.Errorf("expected 'name: worker'")
	}
	if !strings.Contains(s, "type: job") {
		t.Errorf("expected 'type: job'")
	}
	if strings.Contains(s, "image:") {
		t.Error("should not contain 'image:' when empty")
	}
	if strings.Contains(s, "port:") {
		t.Error("should not contain 'port:' when 0")
	}
}

func TestToYAMLVersionOnly(t *testing.T) {
	spec := &Spec{
		Project: Project{Name: "ver", Version: "0.1.0"},
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "version: 0.1.0") {
		t.Errorf("expected 'version: 0.1.0'")
	}
	if strings.Contains(s, "description:") {
		t.Error("should not contain 'description:' when empty")
	}
}

func TestParseServiceWithOnlyImage(t *testing.T) {
	input := []byte(`
service "imgonly" {
  image = "nginx:latest"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	svc := spec.Services["imgonly"]
	if svc.Image != "nginx:latest" {
		t.Errorf("expected 'nginx:latest', got %q", svc.Image)
	}
	if svc.Port != 0 {
		t.Errorf("expected port 0, got %d", svc.Port)
	}
	if svc.Type != "" {
		t.Errorf("expected empty type, got %q", svc.Type)
	}
}

func TestParseInfraEngineVariants(t *testing.T) {
	for _, engine := range []string{"docker", "kubernetes", "podman", "nomad"} {
		input := []byte(`infra "infra" {
  engine = "` + engine + `"
}`)
		spec, err := Parse(input, "test.hcl")
		if err != nil {
			t.Fatal(err)
		}
		if spec.Infra.Engine != engine {
			t.Errorf("expected engine %q, got %q", engine, spec.Infra.Engine)
		}
	}
}

func TestParseAndToYAMLRoundTrip(t *testing.T) {
	input := []byte(`
project "roundtrip" {
  version     = "1.0.0"
  description = "Round trip test"
}

service "api" {
  image = "myapp-api"
  port  = 8080
  type  = "backend"
}

infra "infra" {
  engine = "docker"
}
`)
	spec, err := Parse(input, "test.hcl")
	if err != nil {
		t.Fatal(err)
	}
	out, err := ToYAML(spec)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "name: roundtrip") {
		t.Error("round trip lost project name")
	}
	if !strings.Contains(s, "version: 1.0.0") {
		t.Error("round trip lost version")
	}
	if !strings.Contains(s, "description: Round trip test") {
		t.Error("round trip lost description")
	}
	if !strings.Contains(s, "engine: docker") {
		t.Error("round trip lost engine")
	}
}

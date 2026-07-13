package renderers

import (
	"strings"
	"testing"
	"text/template"
)

func TestRenderSimpleTemplate(t *testing.T) {
	r := NewRenderer()
	result, err := r.Render("hello {{.Name}}", map[string]string{"Name": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(result))
	}
}

func TestRenderNamedTemplate(t *testing.T) {
	r := NewRenderer()
	result, err := r.RenderNamed("greeting", "hi {{.Name}}", map[string]string{"Name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hi test" {
		t.Fatalf("expected 'hi test', got '%s'", string(result))
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	r := NewRenderer()
	_, err := r.Render("hello {{.Name", nil)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRenderTemplateWithFunctions(t *testing.T) {
	funcs := template.FuncMap{
		"upper": func(s string) string {
			result := ""
			for _, c := range s {
				if c >= 'a' && c <= 'z' {
					result += string(c - 32)
				} else {
					result += string(c)
				}
			}
			return result
		},
	}
	result, err := RenderWithFuncs("test", "hello {{upper .Name}}", map[string]string{"Name": "world"}, funcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello WORLD" {
		t.Fatalf("expected 'hello WORLD', got '%s'", string(result))
	}
}

func TestRenderTemplateWithNestedData(t *testing.T) {
	type Data struct {
		Project string
		Modules []string
	}
	data := Data{
		Project: "my-project",
		Modules: []string{"auth", "user", "api"},
	}
	tmpl := "project: {{.Project}}\nmodules: {{range .Modules}}{{.}} {{end}}"
	result, err := RenderTemplate("test", tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "project: my-project\nmodules: auth user api "
	if string(result) != expected {
		t.Fatalf("expected '%s', got '%s'", expected, string(result))
	}
}

func TestTemplateData(t *testing.T) {
	data := TemplateData{
		Project: "acme-api",
		Module:  "auth",
		Port:    8080,
	}
	tmpl := "{{.Project}} {{.Module}} :{{.Port}}"
	result, err := RenderTemplate("test", tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "acme-api auth :8080" {
		t.Fatalf("unexpected result: %s", string(result))
	}
}

func TestRenderEmptyTemplate(t *testing.T) {
	r := NewRenderer()
	result, err := r.Render("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got '%s'", string(result))
	}
}

func TestRenderNilData(t *testing.T) {
	r := NewRenderer()
	result, err := r.Render("static text", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "static text" {
		t.Fatalf("expected 'static text', got '%s'", string(result))
	}
}

func TestRenderWithFuncsInvalidTemplate(t *testing.T) {
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
	}
	_, err := RenderWithFuncs("test", "hello {{.Name", nil, funcs)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRenderWithFuncsExecuteError(t *testing.T) {
	funcs := template.FuncMap{}
	_, err := RenderWithFuncs("test", "hello {{.Name}}", nil, funcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderWithFuncsNilFuncMap(t *testing.T) {
	_, err := RenderWithFuncs("test", "hello {{.Name}}", map[string]string{"Name": "world"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderNamedEmptyTemplate(t *testing.T) {
	r := NewRenderer()
	result, err := r.RenderNamed("empty", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got '%s'", string(result))
	}
}

func TestRenderNamedInvalidTemplate(t *testing.T) {
	r := NewRenderer()
	_, err := r.RenderNamed("bad", "unclosed {{", nil)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRenderTemplateDataFull(t *testing.T) {
	data := TemplateData{
		Project:  "my-proj",
		Module:   "auth",
		Service:  "api",
		Port:     9090,
		Kind:     "grpc",
		Version:  "2.0.0",
		Package:  "com.example",
		Language: "go",
		Attributes: map[string]string{
			"env": "prod",
		},
	}
	tmpl := "{{.Project}}/{{.Module}} {{.Service}}:{{.Port}} {{.Kind}} {{.Version}} {{.Package}} {{.Language}} env={{index .Attributes \"env\"}}"
	result, err := RenderTemplate("test", tmpl, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "my-proj/auth api:9090 grpc 2.0.0 com.example go env=prod"
	if string(result) != expected {
		t.Fatalf("expected '%s', got '%s'", expected, string(result))
	}
}

func TestRenderWithFuncsEmptyFuncMap(t *testing.T) {
	funcs := template.FuncMap{}
	result, err := RenderWithFuncs("test", "hello world", nil, funcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(result))
	}
}

func TestRenderTemplateConditional(t *testing.T) {
	type Data struct {
		Show bool
		Name string
	}
	tmpl := "{{if .Show}}{{.Name}}{{else}}hidden{{end}}"
	result, err := RenderTemplate("test", tmpl, Data{Show: true, Name: "visible"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "visible" {
		t.Fatalf("expected 'visible', got '%s'", string(result))
	}
	result, err = RenderTemplate("test", tmpl, Data{Show: false, Name: "visible"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hidden" {
		t.Fatalf("expected 'hidden', got '%s'", string(result))
	}
}

func TestNewRendererReturnsDefaultRenderer(t *testing.T) {
	r := NewRenderer()
	if _, ok := r.(DefaultRenderer); !ok {
		t.Fatal("expected DefaultRenderer type")
	}
}

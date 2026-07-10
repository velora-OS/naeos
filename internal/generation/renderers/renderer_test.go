package renderers

import (
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

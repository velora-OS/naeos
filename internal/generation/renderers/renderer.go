package renderers

import (
	"bytes"
	"fmt"
	"text/template"
)

type Renderer interface {
	Render(tmpl string, data any) ([]byte, error)
	RenderNamed(name, tmpl string, data any) ([]byte, error)
}

type DefaultRenderer struct{}

func NewRenderer() Renderer {
	return DefaultRenderer{}
}

func (DefaultRenderer) Render(tmpl string, data any) ([]byte, error) {
	return RenderTemplate("default", tmpl, data)
}

func (DefaultRenderer) RenderNamed(name, tmpl string, data any) ([]byte, error) {
	return RenderTemplate(name, tmpl, data)
}

func RenderTemplate(name, tmpl string, data any) ([]byte, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

func RenderWithFuncs(name, tmpl string, data any, funcs template.FuncMap) ([]byte, error) {
	t, err := template.New(name).Funcs(funcs).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}

type TemplateData struct {
	Project    string
	Module     string
	Service    string
	Port       int
	Kind       string
	Version    string
	Package    string
	Language   string
	Attributes map[string]string
}

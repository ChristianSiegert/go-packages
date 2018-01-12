package pages

import (
	"errors"
	"html/template"
)

// Template is a collection of (nested) template files.
type Template struct {
	funcMap  template.FuncMap
	paths    []string
	template *template.Template
}

// NewTemplate creates a template from template files specified by paths. If the
// template files are supposed to use functions other than the built-in Go
// functions, these functions must be provided through funcMap.
func NewTemplate(funcMap template.FuncMap, paths ...string) (*Template, error) {
	if len(paths) == 0 {
		return nil, errors.New("pages: no template path provided")
	}

	tpl, err := load(funcMap, paths...)
	if err != nil {
		return nil, err
	}

	return &Template{
		funcMap:  funcMap,
		paths:    paths,
		template: tpl,
	}, nil
}

// MustNewTemplate calls NewTemplate. It panics on error.
func MustNewTemplate(funcMap template.FuncMap, paths ...string) *Template {
	template, err := NewTemplate(funcMap, paths...)
	if err != nil {
		panic(err)
	}
	return template
}

// Reload parses the template files again.
func (t *Template) Reload() error {
	var err error
	t.template, err = load(t.funcMap, t.paths...)
	return err
}

// load parses all files specified by paths.
func load(funcMap template.FuncMap, paths ...string) (*template.Template, error) {
	return template.New("root").Funcs(funcMap).ParseFiles(paths...)
}

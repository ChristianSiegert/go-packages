package pages

import (
	"errors"
	"html/template"
)

// Template is a wrapper around template.Template. It stores information that
// allows it to reload template files.
type Template struct {
	funcMap  template.FuncMap
	paths    []string
	template *template.Template
}

// Reload parses the template files again.
func (t *Template) Reload() error {
	tpl, err := load(t.funcMap, t.paths...)
	t.template = tpl
	return err
}

// NewTemplate loads a template consisting of template files specified by paths.
// If the template uses functions other than the built-in Go functions, they
// must be provided by funcMap.
func NewTemplate(funcMap template.FuncMap, paths ...string) (*Template, error) {
	if len(paths) == 0 {
		return nil, errors.New("pages.NewTemplate: no template path provided")
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

// MustNewTemplate calls NewTemplate. If NewTemplate returns an error, the
// function panics.
func MustNewTemplate(funcMap template.FuncMap, paths ...string) *Template {
	return Must(NewTemplate(funcMap, paths...))
}

// Must panics if err is not nil.
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// load parses all files specified by paths.
func load(funcMap template.FuncMap, paths ...string) (*template.Template, error) {
	return template.New("root").Funcs(funcMap).ParseFiles(paths...)
}

package pages

import "html/template"

// Template is a wrapper around template.Template. It stores information that
// allows it to reload template files.
type Template struct {
	contentTemplatePath string
	funcMap             template.FuncMap
	rootTemplatePath    string
	template            *template.Template
}

// Reload parses the template files again.
func (t *Template) Reload() error {
	tpl, err := load(t.rootTemplatePath, t.contentTemplatePath, t.funcMap)
	t.template = tpl
	return err
}

// NewTemplate loads a template consisting of two template files.
func NewTemplate(rootTemplatePath, contentTemplatePath string, funcMap template.FuncMap) (*Template, error) {
	tpl, err := load(rootTemplatePath, contentTemplatePath, funcMap)
	if err != nil {
		return nil, err
	}

	return &Template{
		contentTemplatePath: contentTemplatePath,
		funcMap:             funcMap,
		rootTemplatePath:    rootTemplatePath,
		template:            tpl,
	}, nil
}

// MustNewTemplate calls NewTemplate. If the root or content template cannot be
// found, the function panics.
func MustNewTemplate(rootTemplatePath, contentTemplatePath string, funcMap template.FuncMap) *Template {
	return Must(NewTemplate(rootTemplatePath, contentTemplatePath, funcMap))
}

// Must panics if err is not nil.
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// load parses the root and content template.
func load(rootTemplatePath, contentTemplatePath string, funcMap template.FuncMap) (*template.Template, error) {
	paths := make([]string, 0, 2)
	paths = append(paths, rootTemplatePath)

	if contentTemplatePath != "" {
		paths = append(paths, contentTemplatePath)
	}

	return template.New("root").Funcs(funcMap).ParseFiles(paths...)
}

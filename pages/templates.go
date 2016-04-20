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

// NewTemplateWithRoot loads a root template and embeds the content template in
// it. The content template is embedded at the location of
// {{template "content" .}} in the root template.
func NewTemplateWithRoot(rootTemplatePath, contentTemplatePath string, funcMap template.FuncMap) (*Template, error) {
	tpl, err := load(rootTemplatePath, contentTemplatePath, funcMap)
	if err != nil {
		return nil, err
	}

	return &Template{
		contentTemplatePath: contentTemplatePath,
		funcMap:             funcMap,
		rootTemplatePath:    RootTemplatePath,
		template:            tpl,
	}, nil
}

// MustNewTemplateWithRoot calls NewTemplateWithRoot. If the root or content
// template cannot be found, the function panics.
func MustNewTemplateWithRoot(rootTemplatePath, contentTemplatePath string, funcMap template.FuncMap) *Template {
	return Must(NewTemplateWithRoot(rootTemplatePath, contentTemplatePath, funcMap))
}

// NewTemplate loads the default root template specified by RootTemplatePath and
// embeds the content template in it. The content template is embedded at the
// location of {{template "content" .}} in the root template.
func NewTemplate(contentTemplatePath string, funcMap template.FuncMap) (*Template, error) {
	return NewTemplateWithRoot(RootTemplatePath, contentTemplatePath, funcMap)
}

// MustNewTemplate calls NewTemplate. If the root or content template cannot be
// found, the function panics.
func MustNewTemplate(contentTemplatePath string, funcMap template.FuncMap) *Template {
	return Must(NewTemplate(contentTemplatePath, funcMap))
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

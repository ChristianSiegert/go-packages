package sessions

import (
	"html/template"
)

type Flash struct {
	Message string
	Type    int
}

// HtmlUnsafe returns the message as type template.HTML without escaping it.
func (f *Flash) HtmlUnsafe() template.HTML {
	return template.HTML(f.Message)
}

func (f *Flash) String() string {
	return f.Message
}

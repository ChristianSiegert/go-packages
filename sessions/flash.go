package sessions

import "html/template"

// Types of flashes.
const (
	FlashTypeError = iota
	FlashTypeInfo
)

// Flash consists of a message that should be displayed to the user, and a type
// that can be used to style the flash appropriately.
type Flash struct {
	Message string
	Type    int
}

// HtmlUnsafe returns f.Message without escaping it. Useful if f.Message
// contains HTML code that should be rendered by the browser.
func (f *Flash) HtmlUnsafe() template.HTML {
	return template.HTML(f.Message)
}

// String returns f.Message.
func (f *Flash) String() string {
	return f.Message
}

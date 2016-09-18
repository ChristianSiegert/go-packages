package sessions

import "html/template"

// Type to identify flash messages that show the user an error.
const FlashTypeError = "error"

// Flash consists of a message that should be displayed to the user, and a type
// that can be used to style the message appropriately.
type Flash struct {
	Message string
	Type    string
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

package sessions

import "html/template"
import "encoding/json"

// Flash consists of a message that should be displayed to the user, and a type
// that can be used to style the message appropriately.
type Flash interface {
	// HTMLUnsafe returns the message without escaping it. This is useful if
	// the message contains HTML code that should be rendered by the browser.
	HTMLUnsafe() template.HTML

	// Message returns the flash’s message.
	Message() string

	// SetMessage sets the flash’s message.
	SetMessage(string)

	// SetType sets the flash’s type.
	SetType(string)

	// String returns the flash’s message.
	String() string

	// Type returns the flash’s type.
	Type() string

	// Flash supports JSON encoding.
	json.Marshaler
}

// flash is an unexported type that implements the Flash interface.
type flash struct {
	flashType string
	message   string
}

// encodableFlash is used for JSON encoding Flash objects.
type encodableFlash struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
}

// NewFlash returns a new instance of Flash.
func NewFlash(message, flashType string) Flash {
	return &flash{
		flashType: flashType,
		message:   message,
	}
}

func (f *flash) HTMLUnsafe() template.HTML {
	return template.HTML(f.Message())
}

func (f *flash) Message() string {
	return f.message
}

func (f *flash) SetMessage(message string) {
	f.message = message
}

func (f *flash) SetType(flashType string) {
	f.flashType = flashType
}

func (f *flash) String() string {
	if f.Type() != "" {
		return f.Type() + ": " + f.Message()
	}
	return f.Message()
}

func (f *flash) Type() string {
	return f.flashType
}

// MarshalJSON returns the JSON encoding of f.
func (f *flash) MarshalJSON() ([]byte, error) {
	temp := &encodableFlash{
		Message: f.Message(),
		Type:    f.Type(),
	}
	return json.Marshal(temp)
}

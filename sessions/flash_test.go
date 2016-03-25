package sessions

import (
	"html/template"
	"testing"
)

func TestHtml(t *testing.T) {
	flash := &Flash{
		Message: "<b>Foo</b>",
	}

	expected := template.HTML(flash.Message)

	if result := flash.HtmlUnsafe(); result != expected {
		t.Errorf("Returned %q, expected %q.", result, expected)
	}
}

func TestString(t *testing.T) {
	flash := &Flash{
		Message: "foo bar",
		Type:    123,
	}

	expected := flash.Message

	if result := flash.String(); result != expected {
		t.Errorf("String() returned %q, expected %q.", result, expected)
	}
}

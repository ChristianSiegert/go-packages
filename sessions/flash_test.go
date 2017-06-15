package sessions

import (
	"bytes"
	"html/template"
	"testing"
)

func TestFlash_HtmlUnsafe(t *testing.T) {
	flash := NewFlash("<b>Foo</b>", "")
	expected := template.HTML(flash.Message())

	if result := flash.HTMLUnsafe(); result != expected {
		t.Errorf("Returned %q, expected %q.", result, expected)
	}
}

func TestFlash_MarshalJSON(t *testing.T) {
	flash := NewFlash("messageA", "typeA")
	expected := []byte("{\"message\":\"messageA\",\"type\":\"typeA\"}")

	if result, err := flash.MarshalJSON(); err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if bytes.Compare(result, expected) != 0 {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestFlash_Message(t *testing.T) {
	flash := NewFlash("foo", "bar")
	expected := "foo"

	if result := flash.Message(); result != expected {
		t.Errorf("Returned %q, expected %q.", result, expected)
	}
}

func TestFlash_SetMessage(t *testing.T) {
	flash := NewFlash("", "")
	flash.SetMessage("messageA")
	expected := "messageA"

	if result := flash.Message(); result != expected {
		t.Errorf("Expected %q, got %q.", expected, result)
	}
}

func TestFlash_SetType(t *testing.T) {
	flash := NewFlash("", "")
	flash.SetType("typeA")
	expected := "typeA"

	if result := flash.Type(); result != expected {
		t.Errorf("Expected %q, got %q.", expected, result)
	}
}

func TestFlash_String(t *testing.T) {
	tests := []struct {
		flashType string
		message   string
		expected  string
	}{
		{"", "foo bar", "foo bar"},
		{"baz", "foo bar", "baz: foo bar"},
	}

	for _, test := range tests {
		flash := NewFlash(test.message, test.flashType)

		if result := flash.String(); result != test.expected {
			t.Errorf("Returned %q, expected %q.", result, test.expected)
		}
	}
}

func TestFlash_Type(t *testing.T) {
	flash := NewFlash("foo", "bar")
	expected := "bar"

	if result := flash.Type(); result != expected {
		t.Errorf("Returned %q, expected %q.", result, expected)
	}
}

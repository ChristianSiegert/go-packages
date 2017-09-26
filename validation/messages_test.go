package validation

import "testing"

func TestMessages_Error(t *testing.T) {
	m := Messages{
		"item1": "message 1",
		"item2": "message 2",
	}

	expected := "validation.Messages{\"item1\":\"message 1\", \"item2\":\"message 2\"}"
	if result := m.Error(); result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}

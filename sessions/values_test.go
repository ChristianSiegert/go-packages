package sessions

import "testing"
import "reflect"

func TestValues_Get(t *testing.T) {
	values := NewValues()
	values.Set("keyA", "valueA")

	if expected, result := "valueA", values.Get("keyA"); result != expected {
		t.Errorf("Expected %q, got %q.", expected, result)
	}
}

func TestValues_Remove(t *testing.T) {
	values := NewValues()
	values.Set("keyA", "valueA")
	values.Set("keyB", "valueB")
	values.Remove("keyA")

	if expected, result := "", values.Get("keyA"); result != expected {
		t.Errorf("Expected %q, got %q.", expected, result)
	}
}

func TestValues_RemoveAll(t *testing.T) {
	values := NewValues()
	values.Set("keyA", "valueA")
	values.Set("keyB", "valueB")
	values.RemoveAll()
	expected := map[string]string{}

	if result := values.GetAll(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestValues_SetAll(t *testing.T) {
	pairs := map[string]string{
		"keyA": "valueA",
		"keyB": "valueB",
	}

	values := NewValues()
	values.SetAll(pairs)

	if result := values.GetAll(); !reflect.DeepEqual(result, pairs) {
		t.Errorf("Expected %v, got %v", pairs, result)
	}
}

// func TestValuesFromJSON(t *testing.T) {
// 	data := []byte("[{\"message\":\"messageA\",\"type\":\"typeA\"},{\"message\":\"messageB\",\"type\":\"typeB\"}]")
// 	expected := []Flash{
// 		NewFlash("messageA", "typeA"),
// 		NewFlash("messageB", "typeB"),
// 	}

// 	if result, err := FlashFromJSON(data); err != nil {
// 		t.Errorf("Unexpected error: %s", err)
// 	} else if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected %#v, got %#v", expected, result)
// 	}
// }

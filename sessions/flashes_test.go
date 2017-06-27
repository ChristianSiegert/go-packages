package sessions

import (
	"reflect"
	"testing"
)

var (
	flashA = NewFlash("a", "type a")
	flashB = NewFlash("b", "type b")
)

func TestFlashes_Add(t *testing.T) {
	flashes := NewFlashes()
	flashes.Add(flashA, flashB)
	expected := []Flash{flashA, flashB}

	if result := flashes.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlashes_AddNew(t *testing.T) {
	flashes := NewFlashes()
	flashC := flashes.AddNew("c", "type c")
	flashD := flashes.AddNew("d", "type d")
	expected := []Flash{flashC, flashD}

	if result := flashes.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlashes_Remove(t *testing.T) {
	flashes := NewFlashes()
	flashA := flashes.AddNew("a")
	flashB := flashes.AddNew("b")
	flashC := flashes.AddNew("c")
	flashD := flashes.AddNew("d")
	flashes.Remove(flashA, flashC)
	expected := []Flash{flashB, flashD}

	if result := flashes.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlashes_RemoveAll(t *testing.T) {
	flashes := NewFlashes()
	flashes.Add(flashA, flashB)
	flashes.RemoveAll()
	expected := []Flash{}

	if result := flashes.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFlashesFromJSON(t *testing.T) {
	data := []byte("[{\"message\":\"messageA\",\"type\":\"typeA\"},{\"message\":\"messageB\",\"type\":\"typeB\"}]")
	expected := []Flash{
		NewFlash("messageA", "typeA"),
		NewFlash("messageB", "typeB"),
	}

	if result, err := FlashesFromJSON(data); err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %#v, got %#v", expected, result)
	}
}

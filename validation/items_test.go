package validation

import "testing"

func TestItems_Add(t *testing.T) {
	items := New()
	items.Add("name1", 1)
	items.Add("name2", "value2")

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	} else if item1, isPresent := items["name1"]; !isPresent {
		t.Errorf("Expected item 1 to be present.")
	} else if item1.value != 1 {
		t.Errorf("Expected value of item 1 to be %d", 1)
	} else if item2, isPresent := items["name2"]; !isPresent {
		t.Errorf("Expected item 2 to be present.")
	} else if item2.value != "value2" {
		t.Errorf("Expected value of item 2 to be %q", "value2")
	}
}

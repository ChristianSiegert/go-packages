package languages

import (
	"fmt"
	"testing"
)

func TestLanguage_AddMulti(t *testing.T) {
	language := NewLanguage("de", "German")
	language.AddMulti(
		"greeting", "Hallo",
		"farewell", "Tschüss",
	)

	t.Logf("%s", language.translations)
	t.Logf("%d", len(language.translations))

	var tests = []struct {
		translationID string
		expected      string
	}{
		{"greeting", "Hallo"},
		{"farewell", "Tschüss"},
	}

	for _, test := range tests {
		if actual := language.T(test.translationID); actual != test.expected {
			t.Errorf("T(%q): Expected %q, got %q.", test.translationID, test.expected, actual)
		}
	}
}

func TestLanguage_T(t *testing.T) {
	german := NewLanguage("de", "German")
	german.Add("greeting", "Hallo {{.Name}}")

	english := NewLanguage("en", "English")
	english.Add("greeting", "Hello {{.Name}}")

	spanish := NewLanguage("es", "Spanish")
	spanish.Add("greeting", "Hola {{.Name}}")
	spanish.Add("farewell", "Adiós {{.Name}}")

	args := []map[string]interface{}{
		map[string]interface{}{"Name": "Christian"},
	}

	var tests = []struct {
		language      *Language
		fallbacks     []*Language
		translationID string
		args          []map[string]interface{}
		expected      string
	}{
		// Test: key exists in language
		{german, nil, "greeting", args, "Hallo Christian"},
		{english, nil, "greeting", args, "Hello Christian"},
		{spanish, nil, "greeting", args, "Hola Christian"},

		// Test: key is missing in language
		{german, nil, "farewell", nil, "farewell"},

		// Test: key is missing in primary language, exists in fallback language
		{german, []*Language{english, spanish}, "farewell", args, "Adiós Christian"},
	}

	for _, test := range tests {
		test.language.Fallbacks = test.fallbacks

		if actual := test.language.T(test.translationID, test.args...); actual != test.expected {
			t.Errorf(
				"Language %q: T(%q, %s): Expected %q, got %q.",
				test.language.Name(),
				test.translationID,
				test.args,
				test.expected,
				actual,
			)
		}
	}
}

func Example() {
	language := NewLanguage("de", "German")
	language.Add("greeting", "Hallo")

	text := language.T("greeting")
	fmt.Println(text)
	// Output: Hallo
}

func Example_withData() {
	language := NewLanguage("de", "German")
	language.Add("greeting", "Hallo {{.Name}}")

	text := language.T("greeting", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)
	// Output: Hallo Christian
}

func Example_withFallbackLanguages() {
	german := NewLanguage("de", "German")
	german.Add("greeting", "Hallo {{.Name}}")

	english := NewLanguage("en", "English")
	english.Add("greeting", "Hello {{.Name}}")

	spanish := NewLanguage("es", "Spanish")
	spanish.Add("greeting", "Hola {{.Name}}")
	spanish.Add("farewell", "Adiós {{.Name}}")

	german.Fallbacks = []*Language{english, spanish}

	text := german.T("greeting", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)

	text = german.T("farewell", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)
	// Output:
	// Hallo Christian
	// Adiós Christian
}

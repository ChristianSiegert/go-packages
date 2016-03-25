package languages

import (
	"fmt"
	"testing"
)

func TestLanguage_TranslateFunc(t *testing.T) {
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
		language          *Language
		fallbackLanguages []*Language
		translationId     string
		args              []map[string]interface{}
		expected          string
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
		T := test.language.TranslateFunc(test.fallbackLanguages...)
		if actual := T(test.translationId, test.args...); actual != test.expected {
			t.Errorf(
				"Language %q: T(%q, %s): Expected %q, got %q.",
				test.language.Name(),
				test.translationId,
				test.args,
				test.expected,
				actual,
			)
		}
	}
}

func Example() {
	german := NewLanguage("de", "German")
	german.Add("greeting", "Hallo")

	T := german.TranslateFunc()
	text := T("greeting")
	fmt.Println(text)
	// Output: Hallo
}

func Example_withData() {
	german := NewLanguage("de", "German")
	german.Add("greeting", "Hallo {{.Name}}")

	T := german.TranslateFunc()
	text := T("greeting", map[string]interface{}{"Name": "Christian"})
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

	T := german.TranslateFunc(english, spanish)

	text := T("greeting", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)

	text = T("farewell", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)
	// Output:
	// Hallo Christian
	// Adiós Christian
}

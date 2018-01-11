package languages_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"text/template"

	"github.com/ChristianSiegert/go-packages/i18n/languages"
)

func MustTemplate(t *testing.T, name, text string) *template.Template {
	tpl, err := template.New(name).Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	return tpl
}

func TestLanguage_Set(t *testing.T) {
	type args struct {
		translationID string
		translation   interface{}
	}

	tests := []struct {
		args     args
		language *languages.Language
		name     string
		want     *languages.Translation
		wantErr  bool
	}{
		{
			args:     args{"comments", "{{.Count}} Comments"},
			language: languages.NewLanguage("en", "English"),
			name:     "translation is of type string (1)",
			want:     &languages.Translation{Other: MustTemplate(t, "comments", "{{.Count}} Comments")},
			wantErr:  false,
		},
		{
			args:     args{"comments", "{{if }}bad syntax}}"},
			language: languages.NewLanguage("en", "English"),
			name:     "translation is of type string (2)",
			wantErr:  true,
		},
		{
			args:     args{"comments", &languages.Translation{One: MustTemplate(t, "comments", "{{.Count}} Comments")}},
			language: languages.NewLanguage("en", "English"),
			name:     "translation is of type *Translation",
			want:     &languages.Translation{One: MustTemplate(t, "comments", "{{.Count}} Comments")},
			wantErr:  false,
		},
		{
			args:     args{"comments", 123},
			language: languages.NewLanguage("en", "English"),
			name:     "translation is of type int",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.language.Set(test.args.translationID, test.args.translation)

			if (err != nil) != test.wantErr {
				t.Errorf("Language.Set() error = %#v, wantErr %#v", err, test.wantErr)
				return
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Language.Set() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestLanguage_SetMulti(t *testing.T) {
	type args struct {
		translations map[string]interface{}
	}

	tests := []struct {
		args     args
		language *languages.Language
		name     string
		wantErr  bool
	}{
		{
			args:     args{translations: map[string]interface{}{"greeting": "Hallo {{.Name}}"}},
			language: languages.NewLanguage("de-de", "German (Germany)"),
			name:     "translation is of type string",
			wantErr:  false,
		},
		{
			args:     args{translations: map[string]interface{}{"greeting": &languages.Translation{Other: MustTemplate(t, "greeting", "Hallo {{.Name}}")}}},
			language: languages.NewLanguage("de-de", "German (Germany)"),
			name:     "translation is of type *Translation",
			wantErr:  false,
		},
		{
			args:     args{translations: map[string]interface{}{"text": 123}},
			language: languages.NewLanguage("de-de", "German (Germany)"),
			name:     "translation is of int",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.language.SetMulti(test.args.translations); (err != nil) != test.wantErr {
				t.Errorf("Language.SetMulti() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestLanguage_Remove(t *testing.T) {
	type args struct {
		translationIDs []string
	}

	tests := []struct {
		args         args
		language     *languages.Language
		name         string
		translations map[string]interface{}
		want         map[string]*languages.Translation
	}{
		{
			args:     args{translationIDs: []string{"key1", "key3"}},
			language: languages.NewLanguage("en", "English"),
			name:     "Remove key1 and key3",
			translations: map[string]interface{}{
				"key1": "value 1",
				"key2": "value 2",
				"key3": "value 3",
			},
			want: map[string]*languages.Translation{
				"key2": &languages.Translation{Other: MustTemplate(t, "key2", "value 2")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.language.SetMulti(test.translations); err != nil {
				t.Error(err)
			}
			test.language.Remove(test.args.translationIDs...)

			if !reflect.DeepEqual(test.language.Translations, test.want) {
				t.Errorf("got %#v, want %#v", test.language.Translations, test.want)
			}
		})
	}
}

func TestLanguage_T(t *testing.T) {
	german := languages.NewLanguage("de", "German")
	german.Set("greeting", "Hallo {{.Name}}")

	english := languages.NewLanguage("en", "English")
	english.Set("greeting", "Hello {{.Name}}")

	spanish := languages.NewLanguage("es", "Spanish")
	spanish.Set("greeting", "Hola {{.Name}}")
	spanish.Set("farewell", "Adiós {{.Name}}")

	var tests = []struct {
		language      *languages.Language
		fallbacks     []*languages.Language
		translationID string
		data          map[string]interface{}
		want          string
	}{
		// Test: translation ID exists in language
		{german, nil, "greeting", map[string]interface{}{"Name": "Christian"}, "Hallo Christian"},
		{english, nil, "greeting", map[string]interface{}{"Name": "Christian"}, "Hello Christian"},
		{spanish, nil, "greeting", map[string]interface{}{"Name": "Christian"}, "Hola Christian"},

		// Test: translation ID is missing in language
		{german, nil, "farewell", nil, "farewell"},

		// Test: translation ID is missing in primary language, exists in fallback language
		{german, []*languages.Language{english, spanish}, "farewell", map[string]interface{}{"Name": "Christian"}, "Adiós Christian"},
	}

	for i, test := range tests {
		t.Run("Test "+strconv.Itoa(i), func(t *testing.T) {
			test.language.Fallbacks = test.fallbacks

			if got := test.language.T(test.translationID, test.data); got != test.want {
				t.Errorf(
					"Language %q: T(%q, %s): Expected %q, got %q.",
					test.language.Name,
					test.translationID,
					test.data,
					test.want,
					got,
				)
			}
		})
	}
}

func Example() {
	language := languages.NewLanguage("de", "German")
	language.Set("greeting", "Hallo")

	text := language.T("greeting")
	fmt.Println(text)
	// Output: Hallo
}

func Example_withData() {
	language := languages.NewLanguage("de", "German")
	language.Set("greeting", "Hallo {{.Name}}")

	text := language.T("greeting", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)
	// Output: Hallo Christian
}

func Example_withFallbackLanguages() {
	german := languages.NewLanguage("de", "German")
	german.Set("greeting", "Hallo {{.Name}}")

	english := languages.NewLanguage("en", "English")
	english.Set("greeting", "Hello {{.Name}}")

	spanish := languages.NewLanguage("es", "Spanish")
	spanish.Set("greeting", "Hola {{.Name}}")
	spanish.Set("farewell", "Adiós {{.Name}}")

	german.Fallbacks = []*languages.Language{english, spanish}

	text := german.T("greeting", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)

	text = german.T("farewell", map[string]interface{}{"Name": "Christian"})
	fmt.Println(text)

	// Output:
	// Hallo Christian
	// Adiós Christian
}

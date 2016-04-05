// package languages provides a mechanism to create and retrieve translations
// for languages. A language can be backed by fallback languages, so if a
// translation does not exist in the language, its fallback languages are
// checked in the specified order. Translations are pure Go code and thus
// compiled with the program, increasing portability.
package languages

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"golang.org/x/net/context"
)

type Language struct {
	// Language code, e.g. “de”, “en” or “en-US”.
	code string

	// Language name, e.g. “German”.
	name string

	// Translation keys and associated translation text.
	translations map[string]*Translation
}

// TranslateFunc looks up translationId, applies templateData to the translation
// template and returns the result.
type TranslateFunc func(translationId string, templateData ...map[string]interface{}) string

// NewLanguage returns a new instance of Language. Code is the language code,
// e.g. “de”, “en” or “en-US”. Name is the language name, e.g. “German”.
func NewLanguage(code, name string) *Language {
	return &Language{
		code:         code,
		name:         name,
		translations: make(map[string]*Translation),
	}
}

// Add adds a translation to the language. If translationId already exists, the
// existing translation is replaced with the provided one. If the template
// cannot be parsed, the method panics.
func (l *Language) Add(translationId, translation string) *Translation {
	tmpl, err := template.New(translationId).Parse(translation)
	if err != nil {
		panic(fmt.Errorf("languages: Parsing translation template failed: %s\n", err))
	}

	t := &Translation{Other: tmpl}
	l.translations[translationId] = t
	return t
}

// Code returns the language code, e.g. “de”, “en” or “en-US”.
func (l *Language) Code() string {
	return l.code
}

// Get retrieves a translation. If translationId cannot be found, nil is
// returned.
func (l *Language) Get(translationId string) *Translation {
	return l.translations[translationId]
}

// Name returns the language name, e.g. “German”.
func (l *Language) Name() string {
	return l.name
}

// Remove removes a translation from the language.
func (l *Language) Remove(translationId string) {
	delete(l.translations, translationId)
}

// TranslateFunc returns a TranslateFunc that is bound to the language and
// provided fallback languages.
func (l *Language) TranslateFunc(fallbackLanguages ...*Language) TranslateFunc {
	// Create list consisting of primary language followed by fallback languages
	languages := make([]*Language, 1+len(fallbackLanguages))
	languages[0] = l
	for i, language := range fallbackLanguages {
		languages[i+1] = language
	}

	return TranslateFunc(func(translationId string, args ...map[string]interface{}) string {
		var templateData map[string]interface{}

		if len(args) > 0 {
			templateData = args[0]
		}

		// Find translation associated with translationId
		for _, language := range languages {
			translation := language.Get(translationId)
			if translation == nil {
				continue
			}

			var buf bytes.Buffer

			if err := translation.Other.Execute(&buf, templateData); err != nil {
				log.Fatalf("languages: Executing template failed: %s\n", err)
			}
			return buf.String()
		}
		return translationId
	})
}

// Translation allows to store multiple versions of a translation, one for each
// quantity grouped defined in section Plural Rules on
// http://cldr.unicode.org/index/cldr-spec/plural-rules .
type Translation struct {
	Zero,
	One,
	Two,
	Few,
	Many,
	Other *template.Template
}

// contextKey is used for attaching TranslateFunc to context.Context.
type contextKey int

// key is the key used for storing and retrieving TranslateFunc from
// context.Context.
const key contextKey = 0

// NewContext returns a new Context carrying TranslateFunc.
func NewContext(ctx context.Context, translateFunc TranslateFunc) context.Context {
	return context.WithValue(ctx, key, translateFunc)
}

// FromContext extracts TranslateFunc from ctx, if present.
func FromContext(ctx context.Context) (TranslateFunc, bool) {
	translateFunc, ok := ctx.Value(key).(TranslateFunc)
	return translateFunc, ok
}

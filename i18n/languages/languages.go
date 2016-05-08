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
)

type Language struct {
	// Language code, e.g. “de”, “en” or “en-US”.
	code string

	// Fallback languages to check when a translation is missing for this
	// language.
	Fallbacks []*Language

	// Language name, e.g. “German”.
	name string

	// Translation keys and associated translation text.
	translations map[string]*Translation
}

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

// T returns the translation associated with translationId. If the translation
// is missing from l, l.Fallbacks will be checked. If the translation is still
// missing, translationId is returned. Args is optional. The first item of args
// is provided to the translation as data, additional items are ignored.
func (l *Language) T(translationId string, args ...map[string]interface{}) string {
	var templateData map[string]interface{}

	if len(args) > 0 {
		templateData = args[0]
	}

	languages := make([]*Language, 0, 1+len(l.Fallbacks))
	languages = append(languages, l)
	languages = append(languages, l.Fallbacks...)

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
}

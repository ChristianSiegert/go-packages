// Package languages provides a mechanism to create and retrieve translations
// for languages. A language can be backed by fallback languages, so if a
// translation does not exist in the language, its fallback languages are
// checked in the specified order. Translations are pure Go code and thus
// compiled with the program, increasing portability.
package languages

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"
)

// Language is a set of translation IDs and their translation text.
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

// Add adds a translation to the language. If translationID already exists, the
// existing translation is replaced with the provided one. If the template
// cannot be parsed, the method panics.
func (l *Language) Add(translationID, translation string) *Translation {
	tmpl, err := template.New(translationID).Parse(translation)
	if err != nil {
		panic(fmt.Errorf("languages.Add: Parsing translation template failed: %s", err))
	}

	t := &Translation{Other: tmpl}
	l.translations[translationID] = t
	return t
}

// AddMulti adds translations to the language. Args with an even index are used
// as translation IDs and items with an odd index are used as translations.
func (l *Language) AddMulti(args ...string) error {
	if len(args)%2 != 0 {
		return errors.New("languages.AddMulti: Number of translation IDs and translations does not match")
	}

	for i, count := 0, len(args); i < count; i += 2 {
		l.Add(args[i], args[i+1])
	}

	return nil
}

// Code returns the language code, e.g. “de”, “en” or “en-US”.
func (l *Language) Code() string {
	return l.code
}

// Get retrieves a translation. If translationID cannot be found, nil is
// returned.
func (l *Language) Get(translationID string) *Translation {
	return l.translations[translationID]
}

// Name returns the language name, e.g. “German”.
func (l *Language) Name() string {
	return l.name
}

// Remove removes a translation from the language.
func (l *Language) Remove(translationID string) {
	delete(l.translations, translationID)
}

// T returns the translation associated with translationID. If the translation
// is missing from l, l.Fallbacks will be checked. If the translation is still
// missing, translationID is returned. Args is optional. The first item of args
// is provided to the translation as data, additional items are ignored.
func (l *Language) T(translationID string, args ...map[string]interface{}) string {
	var templateData map[string]interface{}

	if len(args) > 0 {
		templateData = args[0]
	}

	languages := make([]*Language, 0, 1+len(l.Fallbacks))
	languages = append(languages, l)
	languages = append(languages, l.Fallbacks...)

	// Find translation associated with translationID
	for _, language := range languages {
		translation := language.Get(translationID)
		if translation == nil {
			continue
		}

		var buf bytes.Buffer

		if err := translation.Other.Execute(&buf, templateData); err != nil {
			log.Fatalf("languages: Executing template failed: %s\n", err)
		}
		return buf.String()
	}
	return translationID
}

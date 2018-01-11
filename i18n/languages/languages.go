// Package languages provides a mechanism to create and retrieve translations
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

// Language is a set of translation IDs and their translation text.
type Language struct {
	// Language code, e.g. “de”, “en” or “en-US”.
	Code string

	// Fallback languages to check when a translation is missing for this
	// language.
	Fallbacks []*Language

	// Language name, e.g. “German”.
	Name string

	// Translation IDs and associated translation.
	Translations map[string]*Translation
}

// NewLanguage returns a new instance of Language. Code is the language code,
// e.g. “de”, “en” or “en-US”. Name is the language name, e.g. “German”.
func NewLanguage(code, name string) *Language {
	return &Language{
		Code:         code,
		Name:         name,
		Translations: make(map[string]*Translation),
	}
}

// Set adds a translation identified by translationID to the language. If a
// translation with the provided translationID already exists, it is replaced.
// translation can be of type string or *Translation.
func (l *Language) Set(translationID string, translation interface{}) (*Translation, error) {
	var t *Translation

	switch translation := translation.(type) {
	case string:
		tpl, err := template.New(translationID).Parse(translation)
		if err != nil {
			return nil, fmt.Errorf("languages: parsing translation failed: %s", err)
		}
		t = &Translation{Other: tpl}
	case *Translation:
		t = translation
	default:
		return nil, fmt.Errorf("languages: unsupported type %T", translation)
	}

	l.Translations[translationID] = t
	return t, nil
}

// SetMulti adds translations to the language. translations is a map of
// translation ID as key and translation as value. If a translation with the
// provided translation ID already exists, it is replaced. translation can be of
// type string or *Translation.
func (l *Language) SetMulti(translations map[string]interface{}) error {
	for translationID, translation := range translations {
		if _, err := l.Set(translationID, translation); err != nil {
			return err
		}
	}
	return nil
}

// Get retrieves a translation. If translationID cannot be found, nil is
// returned.
func (l *Language) Get(translationID string) *Translation {
	return l.Translations[translationID]
}

// Remove removes translations from the language.
func (l *Language) Remove(translationIDs ...string) {
	for _, translationID := range translationIDs {
		delete(l.Translations, translationID)
	}
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

		// TODO: Pick plural group based on quantity.
		if err := translation.Other.Execute(&buf, templateData); err != nil {
			log.Printf("languages: executing template %q with data %#v for language %s %s failed: %s\n", translationID, templateData, l.Code, l.Name, err)
			return translationID
		}
		return buf.String()
	}
	return translationID
}

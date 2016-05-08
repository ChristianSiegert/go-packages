package languages

import "text/template"

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

package languages

import "text/template"

// Translation allows to store multiple versions of a translation, one for each
// plural group. See <http://cldr.unicode.org/index/cldr-spec/plural-rules>.
type Translation struct {
	Zero,
	One,
	Two,
	Few,
	Many,
	Other *template.Template
}

package languages

import "golang.org/x/net/context"

// key is an unexported type for context keys defined in this package. This
// prevents collisions with context keys defined in other packages.
type key int

const keyLanguage key = 0

// NewContext returns a new Context that carries Language.
func NewContext(ctx context.Context, language *Language) context.Context {
	return context.WithValue(ctx, keyLanguage, language)
}

// FromContext returns Language stored in ctx.
func FromContext(ctx context.Context) (*Language, bool) {
	language, ok := ctx.Value(keyLanguage).(*Language)
	return language, ok
}

package languages

import (
	"context"
	"errors"
)

type contextKey int

// key is used to identify the value carried by context.
const key contextKey = 0

// ErrNoLanguage is the error returned by FromContext when context does not
// carry a language.
var ErrNoLanguage = errors.New("language not carried by context")

// NewContext returns a new context that carries language.
func NewContext(ctx context.Context, language *Language) context.Context {
	return context.WithValue(ctx, key, language)
}

// FromContext returns the language carried by ctx. If no language is carried,
// error is ErrNoLanguage.
func FromContext(ctx context.Context) (*Language, error) {
	language, ok := ctx.Value(key).(*Language)
	if !ok {
		return nil, ErrNoLanguage
	}
	return language, nil
}

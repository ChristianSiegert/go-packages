// Package chttp stores *http.Request and http.ResponseWriter in
// context.Context.
package chttp

import (
	"net/http"

	"golang.org/x/net/context"
)

// key is an unexported type for keys defined in this package. This prevents
// collisions with keys defined in other packages.
type key int

const (
	keyRequest key = iota
	keyResponseWriter
)

// NewContext returns a new Context that carries ResponseWriter and Request.
func NewContext(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request) context.Context {
	ctx = context.WithValue(ctx, keyRequest, request)
	ctx = context.WithValue(ctx, keyResponseWriter, responseWriter)
	return ctx
}

// FromContext returns ResponseWriter and Request stored in ctx.
func FromContext(ctx context.Context) (http.ResponseWriter, *http.Request, bool) {
	if ctx == nil {
		return nil, nil, false
	}

	request, ok := ctx.Value(keyRequest).(*http.Request)
	if !ok {
		return nil, nil, false
	}

	responseWriter, ok := ctx.Value(keyResponseWriter).(http.ResponseWriter)
	if !ok {
		return nil, nil, false
	}

	return responseWriter, request, true
}

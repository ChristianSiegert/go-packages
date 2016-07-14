// Package routes provides compact and readable routes that can be used with
// <github.com/julienschmidth/httprouter>.
package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

type Handle func(context.Context, *http.Request, httprouter.Params)

type Route struct {
	Handle  Handle
	Methods []string
	Path    string
}

func New(path string, handle Handle, methods ...string) *Route {
	route := &Route{
		Handle:  handle,
		Methods: methods,
		Path:    path,
	}
	return route
}

// Get makes the route accept GET requests.
func (r *Route) Get() *Route {
	r.Methods = append(r.Methods, http.MethodGet)
	return r
}

// Post makes the route accept POST requests.
func (r *Route) Post() *Route {
	r.Methods = append(r.Methods, http.MethodPost)
	return r
}

// ToHttpRouterHandle makes the route compatible to
// github.com/julienschmidt/httprouter routes.
func (r *Route) ToHttpRouterHandle(f func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params) context.Context) httprouter.Handle {
	return func(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx := context.Background()
		if f != nil {
			ctx = f(ctx, responseWriter, request, params)
			if ctx == nil {
				return
			}
		}
		r.Handle(ctx, request, params)
	}
}

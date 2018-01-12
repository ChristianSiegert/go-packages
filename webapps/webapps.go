// Package webapps handles the creation of routes.
package webapps

import (
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

// Handle responds to an HTTP request.
type Handle func(http.ResponseWriter, *http.Request, httprouter.Params) error

// Middleware envelops Handle to intercept HTTP requests and modify responses.
type Middleware func(Handle) Handle

// WebApp represents a web application or web site. Router gives access to the
// underlying router and its settings. OnError and OnPanic can be overwritten
// by custom functions to handle errors and panics.
type WebApp struct {
	middlewares []Middleware

	// OnError is called after a Handle returned an error.
	OnError func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, err error)

	// OnPanic is called after a Handle panicked.
	OnPanic func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, recoveryInfo interface{})

	// Router is the underlying router.
	Router *httprouter.Router

	serverHost string
	serverPort string
}

// New returns a new WebApp. OnError and OnPanic are initalized with a default
// function for handling errors and panics, and can be overwritten by a custom
// function.
func New(host, port string) *WebApp {
	return &WebApp{
		OnError:    onError,
		OnPanic:    onPanic,
		Router:     httprouter.New(),
		serverHost: host,
		serverPort: port,
	}
}

// Middleware adds a function that is executed before any Handle is executed.
// Middlewares added after calling Route are ignored.
func (w *WebApp) Middleware(middleware Middleware) {
	w.middlewares = append(w.middlewares, middleware)
}

// Route associates a URL path with a Handle.
func (w *WebApp) Route(path string, handle Handle, methods ...string) {
	for _, method := range methods {
		for _, middleware := range w.middlewares {
			handle = middleware(handle)
		}

		h := func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			defer func() {
				if r := recover(); r != nil {
					w.OnPanic(writer, request, params, r)
				}
			}()

			if err := handle(writer, request, params); err != nil {
				w.OnError(writer, request, params, err)
			}
		}

		w.Router.Handle(method, path, h)
	}
}

// Start starts the HTTP server.
func (w *WebApp) Start() error {
	serverAddress := w.serverHost + ":" + w.serverPort
	return http.ListenAndServe(serverAddress, w.Router)
}

// StartWithTLS starts the HTTP server with TLS (Transport Layer Security).
func (w *WebApp) StartWithTLS(certificatePath, keyPath string) error {
	serverAddress := w.serverHost + ":" + w.serverPort
	return http.ListenAndServeTLS(serverAddress, certificatePath, keyPath, w.Router)
}

func onError(writer http.ResponseWriter, request *http.Request, params httprouter.Params, err error) {
	logger.Printf("error %s %s: %s", request.Method, request.URL, err)
	http.Error(writer, "internal server error", http.StatusInternalServerError)
}

func onPanic(writer http.ResponseWriter, request *http.Request, params httprouter.Params, recoveryInfo interface{}) {
	_, file, line, _ := runtime.Caller(4)
	logger.Printf("panic %s %s: %s:%d %+v", request.Method, request.URL, path.Base(file), line, recoveryInfo)
	http.Error(writer, "internal server error", http.StatusInternalServerError)
}

package webapps

import (
	"net/http"

	"github.com/ChristianSiegert/go-packages/i18n/languages"
	"github.com/ChristianSiegert/go-packages/sessions"
	"github.com/julienschmidt/httprouter"
)

// Hook is a function that is called before a route’s handle function.
type Hook func(httprouter.Handle) httprouter.Handle

// WebApp represents a web application or web site.
type WebApp struct {
	hooks      []Hook
	router     *httprouter.Router
	serverHost string
	serverPort string
}

// New returns a new WebApp.
func New(host, port string) *WebApp {
	return &WebApp{
		router:     httprouter.New(),
		serverHost: host,
		serverPort: port,
	}
}

// AddRoute adds a route.
func (w *WebApp) AddRoute(path string, handle httprouter.Handle, methods ...string) {
	for _, method := range methods {
		for _, hook := range w.hooks {
			handle = hook(handle)
		}
		w.router.Handle(method, path, handle)
	}
}

// AddFileDir makes files stored in dirPath accessible at urlPath. urlPath must
// end with “/*filepath”.
func (w *WebApp) AddFileDir(urlPath, dirPath string) {
	w.router.ServeFiles(urlPath, http.Dir(dirPath))
}

// AddHook adds a function that is executed before httprouter.Handle from
// AddRoute executes. Hooks added after calling AddRoute are ignored.
func (w *WebApp) AddHook(hook Hook) {
	w.hooks = append(w.hooks, hook)
}

// SetNotFound sets a NotFound handler the router uses when no route matches.
func (w *WebApp) SetNotFound(handler http.Handler) {
	w.router.NotFound = handler
}

// Start starts the HTTP server.
func (w *WebApp) Start() error {
	serverAddress := w.serverHost + ":" + w.serverPort
	return http.ListenAndServe(serverAddress, w.router)
}

// StartWithTLS starts the HTTP server with TLS (Transport Layer Security).
func (w *WebApp) StartWithTLS(certificatePath, keyPath string) error {
	serverAddress := w.serverHost + ":" + w.serverPort
	return http.ListenAndServeTLS(serverAddress, certificatePath, keyPath, w.router)
}

// LanguageHook returns a hook that adds the user-selected language to the
// request context. param is the name of the route parameter that contains the
// language code. defaultURL is the URL to redirect to when the requested
// language is not supported.
func LanguageHook(param string, langs map[string]*languages.Language, defaultURL string) Hook {
	return func(handle httprouter.Handle) httprouter.Handle {
		return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			languageCode := params.ByName(param)

			language, ok := langs[languageCode]
			if !ok {
				http.Redirect(writer, request, defaultURL, http.StatusSeeOther)
				return
			}

			context := languages.NewContext(request.Context(), language)
			request = request.WithContext(context)

			handle(writer, request, params)
		}
	}
}

// SessionHook returns a hook that adds the session to the request context.
func SessionHook(sessionStore sessions.Store) Hook {
	return func(handle httprouter.Handle) httprouter.Handle {
		return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			if session, err := sessionStore.Get(writer, request); err == nil {
				context := sessions.NewContext(request.Context(), session)
				request = request.WithContext(context)
			}

			handle(writer, request, params)
		}
	}
}

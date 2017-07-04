package webapps

import (
	"log"
	"net/http"

	"github.com/ChristianSiegert/go-packages/i18n/languages"
	"github.com/ChristianSiegert/go-packages/sessions"
	"github.com/julienschmidt/httprouter"
)

// ParamLanguage is the route parameter that stores the language code.
const ParamLanguage = "lang"

// WebApp represents a web application or web site.
type WebApp struct {
	// Default language to redirect to when requested language is not supported.
	defaultLanguageCode string

	languages    map[string]*languages.Language
	logger       *log.Logger
	router       *httprouter.Router
	serverHost   string
	serverPort   string
	sessionStore sessions.Store
}

// New returns a new WebApp.
func New(host, port string, logger *log.Logger, sessionStore sessions.Store) *WebApp {
	return &WebApp{
		languages:    make(map[string]*languages.Language, 1),
		logger:       logger,
		router:       httprouter.New(),
		serverHost:   host,
		serverPort:   port,
		sessionStore: sessionStore,
	}
}

// AddLanguage adds a language to w.
func (w *WebApp) AddLanguage(language *languages.Language, isDefault bool) {
	w.languages[language.Code()] = language
	if isDefault {
		w.defaultLanguageCode = language.Code()
	}
}

// AddRoute adds a route to w.
func (w *WebApp) AddRoute(path string, handle httprouter.Handle, methods ...string) {
	for _, method := range methods {
		handle = w.handleLanguage(handle)
		handle = w.handleSession(handle)
		w.router.Handle(method, path, handle)
	}
}

// AddFileDir makes files stored in dirPath accessible at urlPath. urlPath must
// end with “/*filepath”.
func (w *WebApp) AddFileDir(urlPath, dirPath string) {
	w.router.ServeFiles(urlPath, http.Dir(dirPath))
}

func (w *WebApp) handleLanguage(handle httprouter.Handle) httprouter.Handle {
	if len(w.languages) == 0 {
		return handle
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		languageCode := params.ByName(ParamLanguage)

		// If language is not supported, redirect to default language
		language, ok := w.languages[languageCode]
		if !ok {
			if w.defaultLanguageCode == "" {
				panic("webapps: No default language set.")
			}
			http.Redirect(writer, request, "/"+w.defaultLanguageCode, http.StatusSeeOther)
			return
		}

		// Add language to context
		context := languages.NewContext(request.Context(), language)
		request = request.WithContext(context)

		// Execute given handle
		handle(writer, request, params)
	}
}

func (w *WebApp) handleSession(handle httprouter.Handle) httprouter.Handle {
	if w.sessionStore == nil {
		return handle
	}

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		// Get session for this request
		session, err := w.sessionStore.Get(writer, request)
		if err != nil {
			log.Println("webapps: " + err.Error())
			http.Error(writer, "Interal Server Error", http.StatusInternalServerError)
			return
		}

		// Add session to context
		context := sessions.NewContext(request.Context(), session)
		request = request.WithContext(context)

		// Execute given handle
		handle(writer, request, params)
	}
}

// Start starts the HTTP server.
func (w *WebApp) Start() error {
	serverAddress := w.serverHost + ":" + w.serverPort
	return http.ListenAndServe(serverAddress, w.router)
}

// StartWithTLS starts the HTTP server with TLS (Transport Layer Security).
func (w *WebApp) StartWithTLS(certificatePath, keyPath string) error {
	serverAddress := w.serverHost + w.serverPort
	return http.ListenAndServeTLS(serverAddress, certificatePath, keyPath, w.router)
}

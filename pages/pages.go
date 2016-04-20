// Package pages provides a data structure for a web pages.
package pages

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ChristianSiegert/go-packages/forms"
	"github.com/ChristianSiegert/go-packages/html"
	"github.com/ChristianSiegert/go-packages/i18n/languages"
	"github.com/ChristianSiegert/go-packages/sessions"
	"golang.org/x/net/context"
)

// Whether NewPage and MustNewPage should reload the provided template.
// Reloading templates on each request is useful to see changes without
// recompiling. In production, reloading should be disabled.
var ReloadTemplates = false

// Path to root template.
var RootTemplatePath = "./templates/index.html"

// Templates that are used when Page.ServeEmpty, Error or Page.ServeNotFound is
// called. If a template is nil, only the HTTP status code is set and nothing is
// rendered. To set a different template, set Template[Empty|Error|NotFound]
// from the init function of package main.
var (
	TemplateEmpty    = MustNewTemplate("", nil)
	TemplateError    = MustNewTemplateWithRoot("./templates/error.html", "./templates/500-internal-server-error.html", nil)
	TemplateNotFound = MustNewTemplate("./templates/404-not-found.html", nil)
)

// SignInUrl is the URL to the page that users are redirected to when
// Page.RequireSignIn is called. If a %s placeholder is present in
// SignInUrl.Path, it is replaced by the page’s language code. E.g.
// “/%s/sign-in” becomes “/en/sign-in” if the page’s language code is “en”.
var SignInUrl = &url.URL{
	Path: "/%s/sign-in",
}

// Page represents a web page.
type Page struct {
	Breadcrumbs []*Breadcrumb

	Data map[string]interface{}

	// Form is an instance of *forms.Form bound to the request.
	Form *forms.Form

	LanguageCode string

	// Name is used to highlight the navigation link of the active page.
	Name string

	Request *http.Request

	ResponseWriter http.ResponseWriter

	Session *sessions.Session

	Template *Template

	// Title of the page that templates can use to populate the HTML <title>
	// element.
	Title string

	TranslateFunc languages.TranslateFunc
}

func NewPage(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, languageCode string, tpl *Template) (*Page, error) {
	session, ok := sessions.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("pages.NewPage: sessions.Session is not provided by ctx.")
	}

	translateFunc, ok := languages.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("pages.NewPage: languages.TranslateFunc is not provided by ctx.")
	}

	if ReloadTemplates {
		if err := tpl.Reload(); err != nil {
			return nil, err
		}
	}

	page := &Page{
		Form:           forms.New(request),
		LanguageCode:   languageCode,
		Request:        request,
		ResponseWriter: responseWriter,
		Session:        session,
		Template:       tpl,
		TranslateFunc:  translateFunc,
	}

	return page, nil
}

// MustNewPage calls NewPage and panics on error.
func MustNewPage(ctx context.Context, responseWriter http.ResponseWriter, request *http.Request, languageCode string, tpl *Template) *Page {
	page, err := NewPage(ctx, responseWriter, request, languageCode, tpl)
	if err != nil {
		panic("pages.MustNewPage: " + err.Error())
	}
	return page
}

func (p *Page) AddBreadcrumb(title string, url *url.URL) *Breadcrumb {
	breadcrumb := &Breadcrumb{
		Title: title,
		Url:   url,
	}

	p.Breadcrumbs = append(p.Breadcrumbs, breadcrumb)
	return breadcrumb
}

// RequireSignIn redirects users to the sign-in page specified by SignInUrl.
// If SignInUrl.RawQuery is empty, the query parameters “r” (referrer) and “t”
// (title of the referrer page) are appended. This allows the sign-in page to
// display a message that page <title> is access restricted, and after
// successful authentication, users can be redirected to <referrer>, the page
// they came from.
func (p *Page) RequireSignIn(pageTitle string) {
	u := &url.URL{
		Scheme:   SignInUrl.Scheme,
		Opaque:   SignInUrl.Opaque,
		User:     SignInUrl.User,
		Host:     SignInUrl.Host,
		Path:     fmt.Sprintf(SignInUrl.Path, p.LanguageCode),
		Fragment: SignInUrl.Fragment,
	}

	if SignInUrl.RawQuery == "" {
		query := &url.Values{}
		query.Add("r", p.Request.URL.Path)
		query.Add("t", base64.URLEncoding.EncodeToString([]byte(pageTitle))) // TODO: Sign or encrypt parameter to prevent tempering by users
		u.RawQuery = query.Encode()
	}

	http.Redirect(p.ResponseWriter, p.Request, u.String(), http.StatusSeeOther)
}

// Serve serves the template “index.html” into which it embeds the content
// template specified by page.Template. HTML comments and whitespace are
// stripped. If page.Template is nil, an empty content template is embedded.
func (p *Page) Serve(ctx context.Context) {
	buffer := bytes.NewBuffer([]byte{})

	if p.Template == nil {
		p.Template = TemplateEmpty
	}

	// If still nil
	if p.Template == nil {
		// context := appengine.NewContext(p.Request)
		// context.Errorf("pages.Serve: Content template is nil. Serving blank page.")
		return
	}

	if err := p.Template.template.ExecuteTemplate(buffer, "index.html", p); err != nil {
		// context := appengine.NewContext(p.Request)
		// context.Errorf(err.Error())
		Error(ctx, p.ResponseWriter, p.Request, p.LanguageCode, p.TranslateFunc, err)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(p.ResponseWriter); err != nil {
		// context := appengine.NewContext(p.Request)
		// context.Errorf(err.Error())
		Error(ctx, p.ResponseWriter, p.Request, p.LanguageCode, p.TranslateFunc, err)
	}
}

// ServeEmpty serves the root template without content template.
func (p *Page) ServeEmpty(ctx context.Context) {
	p.Template = TemplateEmpty
	p.Serve(ctx)
}

// ServeNotFound serves a page that tells the user the requested page does not
// exist.
func (page *Page) ServeNotFound(ctx context.Context) {
	page.ResponseWriter.WriteHeader(http.StatusNotFound)
	page.Template = TemplateNotFound
	page.Serve(ctx)
}

// ServeUnauthorized serves a page that tells the user the requested page cannot
// be accessed due to insufficient access rights.
func (p *Page) ServeUnauthorized(ctx context.Context) {
	p.Session.AddFlashErrorMessage(p.TranslateFunc("err_unauthorized_access"))
	p.ResponseWriter.WriteHeader(http.StatusUnauthorized)
	p.ServeEmpty(ctx)
}

// ServeWithError is similar to Serve, but additionally an error flash message
// is displayed to the user saying that an internal problem occurred. Err is not
// displayed but written to the error log. This method is useful if the user
// should be informed of a problem while the state, e.g. a filled in form, is
// preserved.
func (p *Page) ServeWithError(ctx context.Context, err error) {
	// context := appengine.NewContext(p.Request)
	// context.Errorf(err.Error())
	p.Session.AddFlashErrorMessage(p.TranslateFunc("err_internal_server_error"))
	p.Serve(ctx)
}

// Error is an alias for pages.Error.
func (p *Page) Error(ctx context.Context, err error) {
	Error(ctx, p.ResponseWriter, p.Request, p.LanguageCode, p.TranslateFunc, err)
}

// T returns the translation associated with translationId. If p.TranslateFunc
// is nil, translationId is returned.
func (p *Page) T(translationId string, templateData ...map[string]interface{}) string {
	if p.TranslateFunc == nil {
		return translationId
	}
	return p.TranslateFunc(translationId, templateData...)
}

// Error serves an error page with a generic error message. Err is not displayed
// to the user but written to the error log.
func Error(
	ctx context.Context,
	responseWriter http.ResponseWriter,
	request *http.Request,
	languageCode string,
	translateFunc languages.TranslateFunc,
	err error,
) {
	// context := appengine.NewContext(request)
	// context.Errorf(err.Error())
	log.Printf(err.Error())

	if TemplateError == nil {
		// context.Errorf("pages.Error: TemplateError is nil.")
		log.Printf("pages.Error: TemplateError is nil.")
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buffer := bytes.NewBuffer([]byte{})

	errorPage, err2 := NewPage(ctx, responseWriter, request, languageCode, nil)
	if err2 != nil {
		// context.Errorf(err2.Error())
		log.Printf(err2.Error())
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	errorPage.Data = map[string]interface{}{
		"Error":          err,
		"IsDevAppServer": true,
	}

	if err := TemplateError.template.ExecuteTemplate(buffer, "error.html", errorPage); err != nil {
		// context.Errorf("pages.Error: Executing template failed: %s", err)
		log.Printf("pages.Error: Executing template failed: %s", err)
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(responseWriter); err != nil {
		// context.Errorf("pages.Error: Writing template to buffer failed: %s", err)
		log.Printf("pages.Error: Writing template to buffer failed: %s", err)
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
	}
}

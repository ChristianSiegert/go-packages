// Package pages provides a data structure for a web pages.
package pages

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ChristianSiegert/go-packages/chttp"
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

	Language *languages.Language

	// Name of the page. Useful in the root template, e.g. to style the
	// navigation link of the current page.
	Name string

	Request *http.Request

	responseWriter http.ResponseWriter

	Session *sessions.Session

	Template *Template

	// Title of the page that templates can use to populate the HTML <title>
	// element.
	Title string
}

func NewPage(ctx context.Context, tpl *Template) (*Page, error) {
	responseWriter, request, ok := chttp.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("pages.NewPage: http.ResponseWriter and http.Request are not provided by ctx.")
	}

	language, ok := languages.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("pages.NewPage: languages.Language is not provided by ctx.")
	}

	session, ok := sessions.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("pages.NewPage: sessions.Session is not provided by ctx.")
	}

	if ReloadTemplates {
		if err := tpl.Reload(); err != nil {
			return nil, err
		}
	}

	form, err := forms.New(ctx)
	if err != nil {
		return nil, err
	}

	page := &Page{
		Form:           form,
		Language:       language,
		Request:        request,
		responseWriter: responseWriter,
		Session:        session,
		Template:       tpl,
	}

	return page, nil
}

// MustNewPage calls NewPage and panics on error.
func MustNewPage(ctx context.Context, tpl *Template) *Page {
	page, err := NewPage(ctx, tpl)
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

func (p *Page) Redirect(urlStr string, code int) {
	http.Redirect(p.responseWriter, p.Request, urlStr, code)
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
		Path:     fmt.Sprintf(SignInUrl.Path, p.Language.Code()),
		Fragment: SignInUrl.Fragment,
	}

	if SignInUrl.RawQuery == "" {
		query := &url.Values{}
		query.Add("r", p.Request.URL.Path)
		query.Add("t", base64.URLEncoding.EncodeToString([]byte(pageTitle))) // TODO: Sign or encrypt parameter to prevent tempering by users
		u.RawQuery = query.Encode()
	}

	p.Redirect(u.String(), http.StatusSeeOther)
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

	if err := p.Template.template.ExecuteTemplate(buffer, path.Base(p.Template.rootTemplatePath), p); err != nil {
		// context := appengine.NewContext(p.Request)
		// context.Errorf(err.Error())
		Error(ctx, err)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(p.responseWriter); err != nil {
		// context := appengine.NewContext(p.Request)
		// context.Errorf(err.Error())
		Error(ctx, err)
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
	page.responseWriter.WriteHeader(http.StatusNotFound)
	page.Template = TemplateNotFound
	page.Serve(ctx)
}

// ServeUnauthorized serves a page that tells the user the requested page cannot
// be accessed due to insufficient access rights.
func (p *Page) ServeUnauthorized(ctx context.Context) {
	p.Session.AddFlashErrorMessage(p.T("err_unauthorized_access"))
	p.responseWriter.WriteHeader(http.StatusUnauthorized)
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
	p.Session.AddFlashErrorMessage(p.T("err_internal_server_error"))
	p.Serve(ctx)
}

// Error is an alias for pages.Error.
func (p *Page) Error(ctx context.Context, err error) {
	Error(ctx, err)
}

// T returns the translation associated with translationId. If p.Language
// is nil, translationId is returned.
func (p *Page) T(translationId string, templateData ...map[string]interface{}) string {
	if p.Language == nil {
		return translationId
	}
	return p.Language.T(translationId, templateData...)
}

// Error serves an error page with a generic error message. Err is not displayed
// to the user but written to the error log.
func Error(ctx context.Context, err error) {
	// context := appengine.NewContext(request)
	// context.Errorf(err.Error())
	log.Printf(err.Error())

	responseWriter, _, ok := chttp.FromContext(ctx)
	if !ok {
		panic("pages.Error: http.ResponseWriter and http.Request are not provided by ctx.")
	}

	if TemplateError == nil {
		// context.Errorf("pages.Error: TemplateError is nil.")
		log.Printf("pages.Error: TemplateError is nil.")
		http.Error(responseWriter, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buffer := bytes.NewBuffer([]byte{})

	errorPage, err2 := NewPage(ctx, nil)
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

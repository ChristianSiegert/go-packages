// Package pages provides a data structure for a web pages.
package pages

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/ChristianSiegert/go-packages/forms"
	"github.com/ChristianSiegert/go-packages/html"
	"github.com/ChristianSiegert/go-packages/i18n/languages"
	"github.com/ChristianSiegert/go-packages/sessions"
)

// Whether NewPage and MustNewPage should reload templates on every request.
// Reloading templates is useful to see changes without recompiling. In
// production, reloading should be disabled.
var ReloadTemplates = false

// Template that is used when Page.ServeEmpty is called.
var TemplateEmpty *Template

// Template that is used when Page.Error is called.
var TemplateError *Template

// Template that is used when Page.ServeNotFound is called.
var TemplateNotFound *Template

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

	request *http.Request

	Session *sessions.Session

	Template *Template

	// Title of the page that templates can use to populate the HTML <title>
	// element.
	Title string

	writer http.ResponseWriter
}

func NewPage(writer http.ResponseWriter, request *http.Request, tpl *Template) (*Page, error) {
	ctx := request.Context()

	language, ok := languages.FromContext(ctx)
	if !ok {
		return nil, errors.New("pages.NewPage: languages.Language is not provided by context.")
	}

	session, ok := sessions.FromContext(ctx)
	if !ok {
		return nil, errors.New("pages.NewPage: sessions.Session is not provided by context.")
	}

	form, err := forms.New(request)
	if err != nil {
		return nil, err
	}

	page := &Page{
		Data:     make(map[string]interface{}),
		Form:     form,
		Language: language,
		request:  request,
		Session:  session,
		Template: tpl,
		writer:   writer,
	}

	return page, nil
}

// MustNewPage calls NewPage. It panics on error.
func MustNewPage(writer http.ResponseWriter, request *http.Request, tpl *Template) *Page {
	page, err := NewPage(writer, request, tpl)
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

// Redirect redirects the client.
func (p *Page) Redirect(url string, code int) {
	http.Redirect(p.writer, p.request, url, code)
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
		query.Add("r", p.request.URL.Path)
		query.Add("t", base64.URLEncoding.EncodeToString([]byte(pageTitle))) // TODO: Sign or encrypt parameter to prevent tempering by users
		u.RawQuery = query.Encode()
	}

	p.Redirect(u.String(), http.StatusSeeOther)
}

// Error serves an error page with a generic error message. Err is not displayed
// to the user but written to the error log.
func (p *Page) Error(err error) {
	log.Println(err.Error())

	if TemplateError == nil {
		log.Println("pages.Page.Error: No error template provided.")
		http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	buffer := bytes.NewBuffer([]byte{})

	p.Data = map[string]interface{}{
		"Error":          err,
		"IsDevAppServer": true,
	}

	if ReloadTemplates {
		if err := p.Template.Reload(); err != nil {
			log.Printf("pages.Page.Error: Reloading template failed: %s\n", p.Template.paths)
			http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	templateName := path.Base(TemplateError.paths[0])
	if err := TemplateError.template.ExecuteTemplate(buffer, templateName, p); err != nil {
		log.Printf("pages.Page.Error: Executing template failed: %s\n", err)
		http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(p.writer); err != nil {
		log.Printf("pages.Page.Error: Writing template to buffer failed: %s\n", err)
		http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Serve serves the root template specified by RootTemplatePath with the content
// template specified by p.Template. HTML comments and whitespace are stripped.
// If p.Template is nil, an empty content template is embedded.
func (p *Page) Serve() {
	buffer := bytes.NewBuffer([]byte{})

	if p.Template == nil {
		p.Template = TemplateEmpty
	}

	// If still nil
	if p.Template == nil {
		p.Error(errors.New("pages.Page.Serve: No template provided."))
		return
	}

	if ReloadTemplates {
		if err := p.Template.Reload(); err != nil {
			p.Error(err)
		}
	}

	templateName := path.Base(p.Template.paths[0])
	if err := p.Template.template.ExecuteTemplate(buffer, templateName, p); err != nil {
		p.Error(err)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(p.writer); err != nil {
		p.Error(err)
	}
}

// ServeEmpty serves the root template without content template.
func (p *Page) ServeEmpty() {
	p.Template = TemplateEmpty
	p.Serve()
}

// ServeNotFound serves a page that tells the user the requested page does not
// exist.
func (p *Page) ServeNotFound() {
	if TemplateNotFound == nil {
		http.Error(p.writer, p.T("err_404_not_found"), http.StatusNotFound)
		return
	}

	p.writer.WriteHeader(http.StatusNotFound)
	p.Template = TemplateNotFound
	p.Title = p.T("err_404_not_found")
	p.Serve()
}

// ServeUnauthorized serves a page that tells the user the requested page cannot
// be accessed due to insufficient access rights.
func (p *Page) ServeUnauthorized() {
	p.Session.AddFlash(p.T("err_unauthorized_access"), sessions.FlashTypeError)
	p.writer.WriteHeader(http.StatusUnauthorized)
	p.ServeEmpty()
}

// ServeWithError is similar to Serve, but additionally an error flash message
// is displayed to the user saying that an internal problem occurred. Err is not
// displayed but written to the error log. This method is useful if the user
// should be informed of a problem while the state, e.g. a filled in form, is
// preserved.
func (p *Page) ServeWithError(err error) {
	log.Println(err.Error())
	p.Session.AddFlash(p.T("err_505_internal_server_error"), sessions.FlashTypeError)
	p.Serve()
}

// T returns the translation associated with translationId. If p.Language
// is nil, translationId is returned.
func (p *Page) T(translationId string, templateData ...map[string]interface{}) string {
	if p.Language == nil {
		return translationId
	}
	return p.Language.T(translationId, templateData...)
}

func Error(writer http.ResponseWriter, request *http.Request, err error) {
	page := MustNewPage(writer, request, TemplateError)
	page.Error(err)
}

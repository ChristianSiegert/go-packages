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

// FlashTypeError is the flash type used for error messages.
var FlashTypeError = "error"

// Log logs errors. The function can be replaced to use a custom logger.
var Log = func(err error) {
	log.Println(err)
}

// ReloadTemplates is a flag for whether NewPage and MustNewPage should reload
// templates on every request. Reloading templates is useful to see changes
// without recompiling. In production, reloading should be disabled.
var ReloadTemplates = false

// TemplateEmpty is used when Page.ServeEmpty is called.
var TemplateEmpty *Template

// TemplateError is used when Page.Error is called.
var TemplateError *Template

// TemplateNotFound is used when Page.ServeNotFound is called.
var TemplateNotFound *Template

// SignInURL is the URL to the page that users are redirected to when
// Page.RequireSignIn is called. If a %s placeholder is present in
// SignInURL.Path, it is replaced by the page’s language code. E.g.
// “/%s/sign-in” becomes “/en/sign-in” if the page’s language code is “en”.
var SignInURL = &url.URL{
	Path: "/%s/sign-in",
}

// Page represents a web page.
type Page struct {
	// Breadcrumbs manages navigation breadcrumbs.
	Breadcrumbs *Breadcrumbs

	Data map[string]interface{}

	// Form is an instance of *forms.Form bound to the request.
	Form *forms.Form

	Language *languages.Language

	// Name of the page. Useful in the root template, e.g. to style the
	// navigation link of the current page.
	Name string

	request *http.Request

	Session sessions.Session

	Template *Template

	// Title of the page that templates can use to populate the HTML <title>
	// element.
	Title string

	writer http.ResponseWriter
}

// NewPage returns a new page.
func NewPage(writer http.ResponseWriter, request *http.Request, tpl *Template) (*Page, error) {
	ctx := request.Context()

	language, err := languages.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	session, err := sessions.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	form, err := forms.New(request)
	if err != nil {
		return nil, err
	}

	page := &Page{
		Breadcrumbs: &Breadcrumbs{},
		Data:        make(map[string]interface{}),
		Form:        form,
		Language:    language,
		request:     request,
		Session:     session,
		Template:    tpl,
		writer:      writer,
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

// FlashAll returns all flashes, removes them from session and saves the session
// if necessary.
func (p *Page) FlashAll() []sessions.Flash {
	flashes := p.Session.Flashes().GetAll()

	if len(flashes) > 0 {
		p.Session.Flashes().RemoveAll()

		if p.Session.IsStored() {
			if err := p.Session.Save(p.writer); err != nil {
				Log(fmt.Errorf("pages.Page.FlashAll: saving session failed: %s", err))
			}
		}
	}

	return flashes
}

// Redirect redirects the client.
func (p *Page) Redirect(url string, code int) {
	http.Redirect(p.writer, p.request, url, code)
}

// RequireSignIn redirects users to the sign-in page specified by SignInURL.
// If SignInURL.RawQuery is empty, the query parameters “r” (referrer) and “t”
// (title of the referrer page) are appended. This allows the sign-in page to
// display a message that page <title> is access restricted, and after
// successful authentication, users can be redirected to <referrer>, the page
// they came from.
func (p *Page) RequireSignIn(pageTitle string) {
	u := &url.URL{
		Scheme:   SignInURL.Scheme,
		Opaque:   SignInURL.Opaque,
		User:     SignInURL.User,
		Host:     SignInURL.Host,
		Path:     fmt.Sprintf(SignInURL.Path, p.Language.Code()),
		Fragment: SignInURL.Fragment,
	}

	if SignInURL.RawQuery == "" {
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
	Log(err)

	if TemplateError == nil {
		Log(errors.New("pages.Page.Error: no error template provided"))
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
			Log(fmt.Errorf("pages.Page.Error: reloading template failed: %s", p.Template.paths))
			http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	templateName := path.Base(TemplateError.paths[0])
	if err := TemplateError.template.ExecuteTemplate(buffer, templateName, p); err != nil {
		Log(fmt.Errorf("pages.Page.Error: executing template failed: %s", err))
		http.Error(p.writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	b := html.RemoveWhitespace(buffer.Bytes())

	if _, err := bytes.NewBuffer(b).WriteTo(p.writer); err != nil {
		Log(fmt.Errorf("pages.Page.Error: writing template to buffer failed: %s", err))
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
		p.Error(errors.New("pages.Page.Serve: no template provided"))
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
	p.Session.Flashes().AddNew(p.T("err_401_unauthorized"), FlashTypeError)
	p.writer.WriteHeader(http.StatusUnauthorized)
	p.ServeEmpty()
}

// ServeWithError is similar to Serve, but additionally an error flash message
// is displayed to the user saying that an internal problem occurred. Err is not
// displayed but written to the error log. This method is useful if the user
// should be informed of a problem while the state, e.g. a filled in form, is
// preserved.
func (p *Page) ServeWithError(err error) {
	Log(err)
	p.Session.Flashes().AddNew(p.T("err_505_internal_server_error"), FlashTypeError)
	p.Serve()
}

// T returns the translation associated with translationID. If p.Language
// is nil, translationID is returned.
func (p *Page) T(translationID string, templateData ...map[string]interface{}) string {
	if p.Language == nil {
		return translationID
	}
	return p.Language.T(translationID, templateData...)
}

// Error serves a new page using the TemplateError template.
func Error(writer http.ResponseWriter, request *http.Request, err error) {
	page := MustNewPage(writer, request, TemplateError)
	page.Error(err)
}

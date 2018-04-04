// Package pages provides template loading and rendering.
package pages

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/ChristianSiegert/go-packages/forms"
	"github.com/ChristianSiegert/go-packages/html"
	"github.com/ChristianSiegert/go-packages/i18n/languages"
	"github.com/ChristianSiegert/go-packages/sessions"
)

// ReloadTemplates is a flag for whether NewPage should reload templates on
// every request. Reloading templates is useful to see changes without
// recompiling. In production, reloading should be disabled.
var ReloadTemplates = false

// Page represents an HTML page.
type Page struct {
	// BaseURL to prepend to redirect URL.
	BaseURL string

	// Breadcrumbs represent a hierarchical navigation.
	Breadcrumbs *Breadcrumbs

	// Data for populating the template.
	Data map[string]interface{}

	// Form helper for creating HTML input elements in the template.
	Form *forms.Form

	// Language to use for displaying text.
	Language *languages.Language

	// Name of the page. Useful for styling links to the current page
	// differently.
	Name string

	request *http.Request

	// Session associated with the request.
	Session sessions.Session

	// Template to render when calling Serve.
	Template *Template

	// Title of the page. Useful for populating the HTML title element.
	Title string

	writer http.ResponseWriter
}

// NewPage returns a new Page.
func NewPage(writer http.ResponseWriter, request *http.Request, tpl *Template) *Page {
	page := &Page{
		Breadcrumbs: &Breadcrumbs{},
		Data:        make(map[string]interface{}),
		Form:        forms.New(request),
		request:     request,
		Template:    tpl,
		writer:      writer,
	}

	return page
}

// FlashAll returns all flashes, removes them from session and saves the session
// if necessary.
func (p *Page) FlashAll() ([]sessions.Flash, error) {
	if p.Session == nil {
		return nil, errors.New("pages: session is nil")
	}

	flashes := p.Session.Flashes().GetAll()

	if len(flashes) == 0 {
		return flashes, nil
	}

	p.Session.Flashes().RemoveAll()

	if p.Session.IsStored() {
		if err := p.Session.Save(p.writer); err != nil {
			return nil, errors.New("pages: saving session failed: " + err.Error())
		}
	}

	return flashes, nil
}

// Redirect redirects the client to destination, using code as HTTP status code.
// If args is provided, destination is formatted with fmt.Sprintf, to which args
// is passed. destination is automatically prefixed with p.BaseURL.
func (p *Page) Redirect(code int, destination string, args ...interface{}) error {
	if len(args) > 0 {
		destination = fmt.Sprintf(destination, args...)
	}
	http.Redirect(p.writer, p.request, p.BaseURL+destination, code)
	return nil
}

// Serve serves the page.
func (p *Page) Serve() error {
	buffer := bytes.NewBuffer([]byte{})

	if p.Template == nil {
		return errors.New("pages: template is nil")
	}

	if ReloadTemplates {
		if err := p.Template.Reload(); err != nil {
			return err
		}
	}

	templateName := path.Base(p.Template.paths[0])
	if err := p.Template.template.ExecuteTemplate(buffer, templateName, p); err != nil {
		return err
	}

	b := html.RemoveWhitespace(buffer.Bytes())
	_, err := bytes.NewBuffer(b).WriteTo(p.writer)
	return err
}

// T returns the translation associated with translationID. If none is
// associated, it returns translationID.
func (p *Page) T(translationID string, templateData ...map[string]interface{}) string {
	if p.Language == nil {
		return translationID
	}
	return p.Language.T(translationID, templateData...)
}

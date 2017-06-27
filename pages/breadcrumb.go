package pages

import "net/url"

// Breadcrumb represents a navigation breadcrumb.
type Breadcrumb struct {
	Title string
	URL   *url.URL
}

package pages

import "net/url"

// Breadcrumb represents a navigation item.
type Breadcrumb struct {
	Title string
	Url   *url.URL
}

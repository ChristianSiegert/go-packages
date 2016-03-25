package pages

import "net/url"

type Breadcrumb struct {
	Title string
	Url   *url.URL
}

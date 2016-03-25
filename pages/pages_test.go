package pages

import (
	"net/url"
	"testing"
)

func TestRemoveBreadcrumb(t *testing.T) {
	page := &Page{}

	var (
		breadcrumb1 = page.AddBreadcrumb("breadcrumb 1", nil)
		breadcrumb2 = page.AddBreadcrumb("breadcrumb 2", nil)
		breadcrumb3 = page.AddBreadcrumb("breadcrumb 3", nil)
	)

	if breadcrumb1 == nil {
		t.Fatal("breadcrumb1 is unexpectedly nil.")
	}

	if breadcrumb2 == nil {
		t.Fatal("breadcrumb2 is unexpectedly nil.")
	}

	if breadcrumb3 == nil {
		t.Fatal("breadcrumb3 is unexpectedly nil.")
	}

	// err := page.RemoveBreadcrumb(breadcrumb2)
	// if err != nil {
	// 	t.Fatalf("Returned unexpected error: %q", err)
	// }

	if len(page.Breadcrumbs) != 2 ||
		page.Breadcrumbs[0] != breadcrumb1 ||
		page.Breadcrumbs[1] != breadcrumb3 {
		t.Fatalf(".RemoveBreadcrumb(…) didn’t remove the expected breadcrumb (%p). %+v", breadcrumb2, page.Breadcrumbs)
	}
}

func TestPage_Title(t *testing.T) {

	tests := []struct {
		title         string
		breadcrumbs   []*Breadcrumb
		expectedTitle string
	}{
		{
			title:         "",
			breadcrumbs:   nil,
			expectedTitle: "default_home_page_title",
		},
		{
			title:         "Foo",
			breadcrumbs:   nil,
			expectedTitle: "Foo",
		},
		{
			title: "",
			breadcrumbs: []*Breadcrumb{
				&Breadcrumb{
					Title: "Home page",
					Url:   &url.URL{Path: "/"},
				},
			},
			expectedTitle: "default_home_page_title",
		},
		{
			title: "",
			breadcrumbs: []*Breadcrumb{
				&Breadcrumb{
					Title: "Home page",
					Url:   &url.URL{Path: "/"},
				},
				&Breadcrumb{
					Title: "Detail page",
				},
			},
			expectedTitle: "Detail page - Panboard",
		},
	}

	for _, test := range tests {
		page := &Page{
			Breadcrumbs: test.breadcrumbs,
			title:       test.title,
		}
		if result := page.Title(); result != test.expectedTitle {
			t.Errorf("page.Title() returned %q, expected %q.", result, test.expectedTitle)
		}
	}
}

func TestPage_SetTitle(t *testing.T) {
	tests := []struct {
		title         string
		expectedTitle string
	}{
		{
			title:         "",
			expectedTitle: "",
		},
		{
			title:         "foo",
			expectedTitle: "foo",
		},
	}

	page := &Page{}
	for _, test := range tests {
		if page.SetTitle(test.title); page.title != test.expectedTitle {
			t.Errorf("page.SetTitle(%q) set title to %q, expected %q.", test.title, page.title, test.expectedTitle)
		}
	}
}

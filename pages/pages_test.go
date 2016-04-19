package pages

import (
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

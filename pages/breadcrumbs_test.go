package pages

import (
	"net/url"
	"reflect"
	"testing"
)

var (
	breadcrumbA = &Breadcrumb{Title: "a", URL: &url.URL{Path: "/a"}}
	breadcrumbB = &Breadcrumb{Title: "b", URL: nil}
)

func TestBreadcrumbs_Add(t *testing.T) {
	breadcrumbs := &Breadcrumbs{}
	breadcrumbs.Add(breadcrumbA, breadcrumbB)
	expected := []*Breadcrumb{breadcrumbA, breadcrumbB}

	if result := breadcrumbs.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBreadcrumbs_AddNew(t *testing.T) {
	breadcrumbs := &Breadcrumbs{}
	breadcrumbC := breadcrumbs.AddNew("c", &url.URL{Path: "/c"})
	breadcrumbD := breadcrumbs.AddNew("d", &url.URL{Path: "/d"})
	expected := []*Breadcrumb{breadcrumbC, breadcrumbD}

	if result := breadcrumbs.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBreadcrumbs_Remove(t *testing.T) {
	breadcrumbC := &Breadcrumb{Title: "c", URL: &url.URL{Path: "/c"}}
	breadcrumbD := &Breadcrumb{Title: "d", URL: &url.URL{Path: "/d"}}

	breadcrumbs := &Breadcrumbs{
		breadcrumbA,
		breadcrumbB,
		breadcrumbC,
		breadcrumbD,
	}
	breadcrumbs.Remove(breadcrumbA, breadcrumbC)
	expected := []*Breadcrumb{breadcrumbB, breadcrumbD}

	if result := breadcrumbs.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBreadcrumbs_RemoveAll(t *testing.T) {
	breadcrumbs := &Breadcrumbs{}
	breadcrumbs.Add(breadcrumbA, breadcrumbB)
	breadcrumbs.RemoveAll()
	expected := []*Breadcrumb{}

	if result := breadcrumbs.GetAll(); !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

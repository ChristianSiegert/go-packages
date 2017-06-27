package pages

import "net/url"

// Breadcrumbs manages navigation breadcrumbs.
type Breadcrumbs []*Breadcrumb

// Add adds breadcrumbs.
func (b *Breadcrumbs) Add(breadcrumbs ...*Breadcrumb) {
	*b = append(*b, breadcrumbs...)
}

// AddNew creates a new breadcrumb and adds it.
func (b *Breadcrumbs) AddNew(title string, url *url.URL) *Breadcrumb {
	breadcrumb := &Breadcrumb{
		Title: title,
		URL:   url,
	}

	b.Add(breadcrumb)
	return breadcrumb
}

// GetAll returns all breadcrumbs.
func (b *Breadcrumbs) GetAll() []*Breadcrumb {
	return []*Breadcrumb(*b)
}

// Remove removes breadcrumbs.
func (b *Breadcrumbs) Remove(breadcrumbs ...*Breadcrumb) {
	bb := b.GetAll()
	for _, breadcrumb := range breadcrumbs {
		for i := 0; i < len(bb); i++ {
			if bb[i] == breadcrumb {
				bb = append(bb[:i], bb[i+1:]...)
				break
			}
		}
	}
	*b = bb
}

// RemoveAll removes all breadcrumbs.
func (b *Breadcrumbs) RemoveAll() {
	*b = Breadcrumbs{}
}

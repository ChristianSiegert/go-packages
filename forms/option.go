package forms

// Option can be passed to Form.Select to populate a <select> element with
// <option> elements.
type Option struct {
	Label string
	Value string
}

type Options []*Option

func (o Options) Len() int {
	return len(o)
}

func (o Options) Less(i, j int) bool {
	return o[i].Label < o[j].Label
}

func (o Options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Package elements provides a data structure for representing HTML elements.
package elements

import (
	"html"
	"html/template"
	"sort"
	"strings"
)

type Element struct {
	Attributes map[string]string
	Children   []*Element
	HasEndTag  bool
	TagName    string
	Text       string
}

// AddAttributeValue appends value to the named attribute’s value, separated by
// a space. If the attribute does not exist, it is created first. This method is
// chainable.
func (e *Element) AddAttributeValue(name, value string) *Element {
	if e == nil {
		return nil
	} else if e.Attributes == nil {
		e.Attributes = make(map[string]string)
	}

	if _, ok := e.Attributes[name]; ok {
		e.Attributes[name] += " "
	}
	e.Attributes[name] += value

	return e
}

// SetAttributeValue replaces the value of the named attribute with the provided
// one. If the attribute does not exist, it is created first. This method is
// chainable.
func (e *Element) SetAttributeValue(name, value string) *Element {
	if e == nil {
		return nil
	} else if e.Attributes == nil {
		e.Attributes = make(map[string]string)
	}

	e.Attributes[name] = value
	return e
}

// Html is the same as String, but the HTML code can be used in templates.
func (e *Element) Html() template.HTML {
	return template.HTML(e.String())
}

// String returns the element as HTML code.
func (e *Element) String() string {
	if e == nil || e.TagName == "" {
		return ""
	}

	capacity := 2 + len(e.Children)

	if len(e.Attributes) > 0 {
		capacity += 1
	}

	if e.HasEndTag {
		capacity += 1
	}

	pieces := make([]string, 0, capacity)
	pieces = append(pieces, "<"+e.TagName)

	if len(e.Attributes) > 0 {
		attributes := make(sort.StringSlice, 0, len(e.Attributes))
		for k, v := range e.Attributes {
			if v == "" {
				attributes = append(attributes, k)
			} else {
				attributes = append(attributes, k+`="`+html.EscapeString(v)+`"`)
			}
		}
		sort.Sort(attributes)
		pieces = append(pieces, " "+strings.Join(attributes, " "))
	}

	pieces = append(pieces, ">")

	for _, child := range e.Children {
		pieces = append(pieces, child.String())
	}

	if e.HasEndTag {
		pieces = append(pieces, html.EscapeString(e.Text)+"</"+e.TagName+">")
	}

	return strings.Join(pieces, "")
}

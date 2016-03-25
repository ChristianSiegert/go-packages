package elements

import (
	"reflect"
	"testing"
)

func TestElement_AddAttributeValue(t *testing.T) {
	tests := []struct {
		element    *Element
		attributes map[string]string
		expected   *Element
	}{
		{
			element: nil,
			attributes: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			},
			expected: nil,
		},
		{
			element: &Element{
				TagName: "div",
			},
			attributes: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			},
			expected: &Element{
				Attributes: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
				TagName: "div",
			},
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"foo1": "bar1",
				},
				TagName: "div",
			},
			attributes: map[string]string{
				"foo1": "bar3",
				"foo2": "bar2",
			},
			expected: &Element{
				Attributes: map[string]string{
					"foo1": "bar1 bar3",
					"foo2": "bar2",
				},
				TagName: "div",
			},
		},
	}

	for i, test := range tests {
		for key, value := range test.attributes {
			test.element.AddAttributeValue(key, value)
		}
		if !reflect.DeepEqual(test.element, test.expected) {
			t.Errorf("Test %d: Element is\n%s\nexpected\n%s\n", i+1, test.element, test.expected)
		}
	}
}

func TestElement_SetAttributeValue(t *testing.T) {
	tests := []struct {
		element    *Element
		attributes map[string]string
		expected   *Element
	}{
		{
			element: nil,
			attributes: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			},
			expected: nil,
		},
		{
			element: &Element{
				TagName: "div",
			},
			attributes: map[string]string{
				"foo1": "bar1",
				"foo2": "bar2",
			},
			expected: &Element{
				Attributes: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
				TagName: "div",
			},
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"foo1": "bar1",
				},
				TagName: "div",
			},
			attributes: map[string]string{
				"foo1": "bar3",
				"foo2": "bar2",
			},
			expected: &Element{
				Attributes: map[string]string{
					"foo1": "bar3",
					"foo2": "bar2",
				},
				TagName: "div",
			},
		},
	}

	for i, test := range tests {
		for key, value := range test.attributes {
			test.element.SetAttributeValue(key, value)
		}
		if !reflect.DeepEqual(test.element, test.expected) {
			t.Errorf("Test %d: Element is\n%s\nexpected\n%s\n", i+1, test.element, test.expected)
		}
	}
}

func TestElement_String(t *testing.T) {
	tests := []struct {
		element  *Element
		expected string
	}{
		// Nil
		{
			element:  nil,
			expected: "",
		},
		// Without tag name
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  false,
				TagName:    "",
				Text:       "",
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  false,
				TagName:    "",
				Text:       `Five special HTML characters < > & ' "`,
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  true,
				TagName:    "",
				Text:       "",
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  true,
				TagName:    "",
				Text:       `Five special HTML characters < > & ' "`,
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: false,
				TagName:   "",
				Text:      "",
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: false,
				TagName:   "",
				Text:      `Five special HTML characters < > & ' "`,
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: true,
				TagName:   "",
				Text:      "",
			},
			expected: "",
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: true,
				TagName:   "",
				Text:      `Five special HTML characters < > & ' "`,
			},
			expected: "",
		},
		// With tag name
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  false,
				TagName:    "foo",
				Text:       "",
			},
			expected: "<foo>",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  false,
				TagName:    "foo",
				Text:       `Five special HTML characters < > & ' "`,
			},
			expected: "<foo>",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  true,
				TagName:    "foo",
				Text:       "",
			},
			expected: "<foo></foo>",
		},
		{
			element: &Element{
				Attributes: nil,
				HasEndTag:  true,
				TagName:    "foo",
				Text:       `Five special HTML characters < > & ' "`,
			},
			expected: "<foo>Five special HTML characters &lt; &gt; &amp; &#39; &#34;</foo>",
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: false,
				TagName:   "foo",
				Text:      "",
			},
			expected: `<foo a="Attribute a value" b c="Attribute c value" d="Attribute d value &lt; &gt; &amp; &#39; &#34;">`,
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: false,
				TagName:   "foo",
				Text:      `Five special HTML characters < > & ' "`,
			},
			expected: `<foo a="Attribute a value" b c="Attribute c value" d="Attribute d value &lt; &gt; &amp; &#39; &#34;">`,
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: true,
				TagName:   "foo",
				Text:      "",
			},
			expected: `<foo a="Attribute a value" b c="Attribute c value" d="Attribute d value &lt; &gt; &amp; &#39; &#34;"></foo>`,
		},
		{
			element: &Element{
				Attributes: map[string]string{
					"d": `Attribute d value < > & ' "`,
					"a": "Attribute a value",
					"c": "Attribute c value",
					"b": "",
				},
				HasEndTag: true,
				TagName:   "foo",
				Text:      `Five special HTML characters < > & ' "`,
			},
			expected: `<foo a="Attribute a value" b c="Attribute c value" d="Attribute d value &lt; &gt; &amp; &#39; &#34;">Five special HTML characters &lt; &gt; &amp; &#39; &#34;</foo>`,
		},
	}

	for i, test := range tests {
		if result := test.element.String(); result != test.expected {
			t.Errorf("Test %d returned\n%s\nexpected\n%s", i+1, result, test.expected)
		}
	}
}

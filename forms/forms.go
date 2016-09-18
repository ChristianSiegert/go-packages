// Package forms provides support for HTML forms and its input elements.
package forms

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ChristianSiegert/go-packages/html/elements"
	"github.com/ChristianSiegert/go-packages/validation"
)

type Form struct {
	request *http.Request

	// ValidationFields is a map of field names and their corresponding
	// *validation.Field. Used to get information about the fields’ validation
	// rules.
	ValidationFields validation.Fields

	// ValidationMessages contains a validation error message for each form
	// field that has an invalid value.
	ValidationMessages validation.Messages
}

func New(request *http.Request) (*Form, error) {
	return &Form{
		request: request,
	}, nil
}

// Error returns a <div class="validation-error"> element that contains the
// validation error message as text.
func (f *Form) Error(fieldName string) *elements.Element {
	validationMessage, ok := f.ValidationMessages[fieldName]
	if !ok || validationMessage == "" {
		return nil
	}

	element := &elements.Element{
		Attributes: map[string]string{
			"class": "validation-error",
		},
		HasEndTag: true,
		TagName:   "div",
		Text:      validationMessage,
	}
	return element
}

// HasError returns whether the field’s value is invalid.
func (f *Form) HasError(fieldName string) bool {
	_, ok := f.ValidationMessages[fieldName]
	return ok
}

// Input returns an <input> element.
func (f *Form) Input(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := &elements.Element{
		Attributes: map[string]string{
			"id":   fieldName,
			"name": fieldName,
		},
		TagName: "input",
	}

	if f.HasError(fieldName) {
		element.Attributes["class"] = "error"
	}

	if placeholder != "" {
		element.Attributes["placeholder"] = placeholder
	}

	if value := f.request.FormValue(fieldName); value != "" {
		element.Attributes["value"] = strings.TrimSpace(value)
	}

	for i, length := 0, len(attributes); i < length; i += 2 {
		if i+1 < length {
			element.AddAttributeValue(attributes[i], attributes[i+1])
		} else {
			element.AddAttributeValue(attributes[i], "")
		}
	}

	if f.ValidationFields == nil {
		return element
	}

	if field, ok := f.ValidationFields[fieldName]; ok {
		for _, rule := range field.Rules {
			if rule.Type == validation.RuleTypeRequired {
				element.Attributes["required"] = ""
			} else if rule.Type == validation.RuleTypeMaxLength {
				if maxLength, ok := rule.Args[0].(int); ok && maxLength > 0 {
					element.Attributes["maxlength"] = strconv.FormatUint(uint64(maxLength), 10)
				}
			} else if rule.Type == validation.RuleTypeMinLength {
				if minLength, ok := rule.Args[0].(int); ok && minLength > 0 {
					element.Attributes["minlength"] = strconv.FormatUint(uint64(minLength), 10)
				}
			}
		}
	}

	return element
}

// Checkbox returns an <input type="checkbox"> element.
func (f *Form) Checkbox(fieldName, value string) *elements.Element {
	element := f.Input(fieldName, "")
	element.Attributes["id"] = fieldName + "-" + value
	element.Attributes["type"] = "checkbox"
	element.Attributes["value"] = value

	if values, ok := f.request.Form[fieldName]; ok {
		for _, v := range values {
			if v == value {
				element.Attributes["checked"] = ""
				break
			}
		}
	}

	return element
}

// Email returns an <input type="email"> element.
func (f *Form) Email(fieldName, placeholder string) *elements.Element {
	element := f.Input(fieldName, placeholder)
	element.Attributes["type"] = "email"

	if _, ok := element.Attributes["maxlength"]; !ok {
		element.Attributes["maxlength"] = "254"
	}
	return element
}

// Number returns an <input type="number"> element.
func (f *Form) Number(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "number"
	return element
}

// Radio returns an <input type="radio"> element.
func (f *Form) Radio(fieldName, value string) *elements.Element {
	element := f.Input(fieldName, "")
	element.Attributes["id"] = fieldName + "-" + value
	element.Attributes["type"] = "radio"
	element.Attributes["value"] = value

	if f.request.FormValue(fieldName) == value {
		element.Attributes["checked"] = ""
	}
	return element
}

// Password returns an <input type="password"> element.
func (f *Form) Password(fieldName, placeholder string) *elements.Element {
	element := f.Input(fieldName, placeholder)
	element.Attributes["type"] = "password"
	return element
}

// Search returns an <input type="search"> element.
func (f *Form) Search(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "search"
	return element
}

// Text returns an <input type="text"> element.
func (f *Form) Text(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "text"
	return element
}

// Select returns a <select> element.
func (f *Form) Select(fieldName string, options []*Option) *elements.Element {
	element := &elements.Element{
		Attributes: map[string]string{
			"id":   fieldName,
			"name": fieldName,
		},
		HasEndTag: true,
		TagName:   "select",
	}

	if len(options) > 0 {
		element.Children = make([]*elements.Element, 0, len(options))

		for _, option := range options {
			optionElement := f.Option(option.Value, option.Label)
			element.Children = append(element.Children, optionElement)

			if postedValue := f.request.FormValue(fieldName); postedValue == option.Value {
				optionElement.SetAttributeValue("selected", "")
			}
		}
	}

	if f.HasError(fieldName) {
		element.Attributes["class"] = "error"
	}

	return element
}

// Option returns an <option> element.
func (f *Form) Option(value, label string) *elements.Element {
	return &elements.Element{
		Attributes: map[string]string{
			"value": value,
		},
		HasEndTag: true,
		TagName:   "option",
		Text:      label,
	}
}

// Textarea returns a <textarea> element.
func (f *Form) Textarea(fieldName, placeholder string) *elements.Element {
	element := &elements.Element{
		Attributes: map[string]string{
			"id":   fieldName,
			"name": fieldName,
		},
		HasEndTag: true,
		TagName:   "textarea",
	}

	if f.HasError(fieldName) {
		element.Attributes["class"] = "error"
	}

	if placeholder != "" {
		element.Attributes["placeholder"] = placeholder
	}

	if value := f.request.FormValue(fieldName); value != "" {
		element.Text = strings.TrimSpace(value)
	}

	if f.ValidationFields == nil {
		return element
	}

	if field, ok := f.ValidationFields[fieldName]; ok {
		for _, rule := range field.Rules {
			if rule.Type == validation.RuleTypeRequired {
				element.Attributes["required"] = ""
			} else if rule.Type == validation.RuleTypeMaxLength {
				if maxLength, ok := rule.Args[0].(int); ok && maxLength > 0 {
					element.Attributes["maxlength"] = strconv.FormatUint(uint64(maxLength), 10)
				}
			} else if rule.Type == validation.RuleTypeMinLength {
				if minLength, ok := rule.Args[0].(int); ok && minLength > 0 {
					element.Attributes["minlength"] = strconv.FormatUint(uint64(minLength), 10)
				}
			}
		}
	}

	return element
}

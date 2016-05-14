// Package forms provides support for HTML forms and its input elements.
package forms

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ChristianSiegert/go-packages/chttp"
	"github.com/ChristianSiegert/go-packages/html/elements"
	"github.com/ChristianSiegert/go-packages/validation"
	"golang.org/x/net/context"
)

type Form struct {
	// ErrorMessages contains error messages that are displayed to the user if
	// the corresponding form fields did not validate successfully. Each form
	// field can only have one error message. The key is the field name, the
	// value is the error message.
	ErrorMessages map[string]string

	request *http.Request

	ValidationFields validation.Fields
}

func New(ctx context.Context) (*Form, error) {
	_, request, ok := chttp.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("forms.New: http.Request is not provided by ctx.")
	}

	return &Form{
		request: request,
	}, nil
}

func (f *Form) Error(fieldName string) *elements.Element {
	errorMessage, ok := f.ErrorMessages[fieldName]
	if !ok || errorMessage == "" {
		return nil
	}

	element := &elements.Element{
		Attributes: map[string]string{
			"class": "validation-error",
		},
		HasEndTag: true,
		TagName:   "div",
		Text:      errorMessage,
	}
	return element
}

func (f *Form) HasError(fieldName string) bool {
	_, ok := f.ErrorMessages[fieldName]
	return ok
}

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

func (f *Form) Email(fieldName, placeholder string) *elements.Element {
	element := f.Input(fieldName, placeholder)
	element.Attributes["type"] = "email"

	if _, ok := element.Attributes["maxlength"]; !ok {
		element.Attributes["maxlength"] = "254"
	}
	return element
}

func (f *Form) Number(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "number"
	return element
}

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

func (f *Form) Password(fieldName, placeholder string) *elements.Element {
	element := f.Input(fieldName, placeholder)
	element.Attributes["type"] = "password"
	return element
}

func (f *Form) Search(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "search"
	return element
}

func (f *Form) Text(fieldName, placeholder string, attributes ...string) *elements.Element {
	element := f.Input(fieldName, placeholder, attributes...)
	element.Attributes["type"] = "text"
	return element
}

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

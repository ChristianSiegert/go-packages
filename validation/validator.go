// Package validation provides means to validate arbitrary data types.
package validation

import (
	"errors"
	"regexp"
)

// Common types of rules. Rule types themselves do nothing, they simply express
// what a rule is supposed to check, so form input fields can be more specific,
// e.g. set the HTML “maxlength” attribute if the rule is of type
// RuleTypeMaxLength.
const (
	RuleTypeEmailAddress = iota + 1
	RuleTypeMaxLength
	RuleTypeMinLength
	RuleTypeRequired
)

// Regular expression for validating an e-mail address
var eMailAddressRegExp = regexp.MustCompile("^[^@]+@[^@]+$")

var ErrValidatorNil = errors.New("validation: validator is nil")

// ErrorMessages is a map whose key is the name of the field that was validated,
// and whose value is an error message intended for display to the user.
type ErrorMessages map[string]string

type Validator struct {
	Fields map[string]*Field
}

// AddField adds a value and its identifying name to validator so it can be
// validated.
func (v *Validator) AddField(name string, value interface{}) *Field {
	field := &Field{
		value: value,
	}

	if v.Fields == nil {
		v.Fields = make(map[string]*Field)
	}

	v.Fields[name] = field
	return field
}

// Validate validates all fields.
func (v *Validator) Validate() (ErrorMessages, error) {
	if v == nil {
		return nil, ErrValidatorNil
	}

	var errorMessages ErrorMessages

	for fieldName, field := range v.Fields {
		if isValid, errorMessage, err := field.Validate(); err != nil {
			return nil, err
		} else if !isValid {
			if errorMessages == nil {
				errorMessages = make(ErrorMessages)
			}
			errorMessages[fieldName] = errorMessage
		}
	}

	return errorMessages, nil
}

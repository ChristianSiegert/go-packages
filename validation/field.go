package validation

import (
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"
)

// Common types of rules. Rule types express what a rule is supposed to check.
// This information can be used to improve input fields, e.g. HTML form fields
// can use attributes that correspond with the rule type.
const (
	RuleTypeEmailAddress = iota + 1
	RuleTypeMaxLength
	RuleTypeMinLength
	RuleTypeRequired
)

// Regular expression for validating an e-mail address.
var eMailAddressRegExp = regexp.MustCompile("^[^@]+@[^@]+$")

// A field can have zero or more validation rules that are used to validate the
// field’s value.
type Field struct {
	Rules []*Rule
	value interface{}
}

// EmailAddress checks if the field’s value is an e-mail address. It only
// checks the length and whether there is exactly one at-sign preceded and
// followed by at least one character.
func (f *Field) EmailAddress(errorMessage string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				return utf8.RuneCountInString(value) <= 254 && eMailAddressRegExp.MatchString(value), nil
			}
			return false, fmt.Errorf("validation.Field.EmailAddress: unsupported value type %T", value)
		},
		Message: errorMessage,
		Type:    RuleTypeEmailAddress,
	})
	return f
}

// Equals checks if the field’s value equals value2.
func (f *Field) Equals(value2 interface{}, message string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case int:
				return value == value2, nil
			case string:
				if value2, ok := value2.(string); !ok || value != value2 {
					return false, nil
				}
				return true, nil
			}
			return false, fmt.Errorf("validation.Field.Equals: unsupported value type %T", value)
		},
		Message: message,
	})
	return f
}

func (f *Field) Func(fn func(value interface{}) (bool, error), message string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func:    fn,
		Message: message,
	})
	return f
}

// MaxLength checks if the field’s value has a maximum length of maxLength.
func (f *Field) MaxLength(maxLength int, message string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				if length := utf8.RuneCountInString(value); length > maxLength {
					return false, nil
				}
				return true, nil
			}
			return false, fmt.Errorf("validation.Field.MaxLength: unsupported value type %T", value)
		},
		Args:    []interface{}{maxLength},
		Message: message,
		Type:    RuleTypeMaxLength,
	})
	return f
}

// MinLength checks if the field’s value has a minimum length of minLength.
func (f *Field) MinLength(minLength int, message string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				if length := utf8.RuneCountInString(value); length < minLength {
					return false, nil
				}
				return true, nil
			}
			return false, fmt.Errorf("validation.Field.MinLength: unsupported value type %T", value)
		},
		Args:    []interface{}{minLength},
		Message: message,
		Type:    RuleTypeMinLength,
	})
	return f
}

// Required checks if the field’s value is non-zero.
func (f *Field) Required(message string) *Field {
	f.Rules = append(f.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case int8:
				return value != 0, nil
			case string:
				return len(value) > 0, nil
			case time.Time:
				return !value.IsZero(), nil
			}
			return false, fmt.Errorf("validation.Field.Required: unsupported value type %T", value)
		},
		Args:    []interface{}{RuleTypeRequired},
		Message: message,
		Type:    RuleTypeRequired,
	})
	return f
}

// Validate checks if the field’s value is valid according to the specified
// validation rules. If it is valid, the function returns true. If it is not
// valid, the rule’s validation error message is returned. If an error
// occurred, the error is returned. A returned error does not mean the value is
// invalid, it solely means something went wrong.
func (f *Field) Validate() (bool, string, error) {
	for _, rule := range f.Rules {
		if isValid, err := rule.Func(f.value); err != nil {
			return false, "", err
		} else if !isValid {
			return false, rule.Message, nil
		}
	}
	return true, "", nil
}

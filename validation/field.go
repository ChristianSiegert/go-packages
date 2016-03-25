package validation

import (
	"fmt"
	"time"
	"unicode/utf8"
)

type Field struct {
	Rules []*Rule
	value interface{}
}

// Required checks if field’s value is an e-mail address. It only checks the
// length and whether there is exactly one at-sign preceded and followed by at
// least one character.
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

// Equals checks if f’s value equals value2.
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

// MaxLength checks if field’s value has a maximum length of maxLength.
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

// MinLength checks if field’s value has a minimum length of minLength.
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

// Required checks if field’s value is non-zero.
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

// Validate checks if field is valid according to the specified rules. If it is
// valid, the function returns true. If field is not valid, an error message is
// returned. If an error occurred, error is returned. A returned error does not
// mean field is invalid, it solely means something went wrong.
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

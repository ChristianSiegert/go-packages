package validation

import (
	"fmt"
	"regexp"
	"strconv"
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

// Item can have zero or more validation rules that are used to validate the
// item’s value.
type Item struct {
	Rules []*Rule
	value interface{}
}

// EmailAddress checks if the item’s value is an e-mail address. It only checks
// the length and whether there is exactly one “at” sign preceded and followed
// by at least one character.
func (i *Item) EmailAddress(message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				return utf8.RuneCountInString(value) <= 254 && eMailAddressRegExp.MatchString(value), nil
			}
			return false, fmt.Errorf("validation.Item.EmailAddress: unsupported value type %T", value)
		},
		Message: message,
		Type:    RuleTypeEmailAddress,
	})
	return i
}

// Equals checks if the item’s value equals value2.
func (i *Item) Equals(value2 interface{}, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
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
			return false, fmt.Errorf("validation.Item.Equals: unsupported value type %T", value)
		},
		Message: message,
	})
	return i
}

func (i *Item) Func(fn func(value interface{}) (bool, error), message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func:    fn,
		Message: message,
	})
	return i
}

// Max checks if the item’s value is equal or less than max.
func (i *Item) Max(max float64, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch v := value.(type) {
			case string:
				number, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return false, err
				}
				return number <= max, nil
			}
			return false, fmt.Errorf("validation.Item.Max: unsupported value type %T", value)
		},
		Args:    []interface{}{max},
		Message: message,
	})
	return i
}

// MaxLength checks if the item’s value has a maximum length of maxLength.
func (i *Item) MaxLength(maxLength int, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				if length := utf8.RuneCountInString(value); length > maxLength {
					return false, nil
				}
				return true, nil
			}
			return false, fmt.Errorf("validation.Item.MaxLength: unsupported value type %T", value)
		},
		Args:    []interface{}{maxLength},
		Message: message,
		Type:    RuleTypeMaxLength,
	})
	return i
}

// Min checks if the item’s value is equal or greater than min.
func (i *Item) Min(min float64, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch v := value.(type) {
			case string:
				number, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return false, err
				}
				return number >= min, nil
			}
			return false, fmt.Errorf("validation.Item.Min: unsupported value type %T", value)
		},
		Args:    []interface{}{min},
		Message: message,
	})
	return i
}

// MinLength checks if the item’s value has a minimum length of minLength.
func (i *Item) MinLength(minLength int, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				if length := utf8.RuneCountInString(value); length < minLength {
					return false, nil
				}
				return true, nil
			}
			return false, fmt.Errorf("validation.Item.MinLength: unsupported value type %T", value)
		},
		Args:    []interface{}{minLength},
		Message: message,
		Type:    RuleTypeMinLength,
	})
	return i
}

// Number checks if the item’s value is a number.
func (i *Item) Number(message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch v := value.(type) {
			case string:
				_, err := strconv.ParseFloat(v, 64)
				return err == nil, nil
			}
			return false, fmt.Errorf("validation.Item.Max: unsupported value type %T", value)
		},
		Message: message,
	})
	return i
}

// Pattern checks if the item’s value matches the regular expression.
func (i *Item) Pattern(pattern *regexp.Regexp, message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case string:
				return pattern.MatchString(value), nil
			}
			return false, fmt.Errorf("validation.Item.Pattern: unsupported value type %T", value)
		},
		Message: message,
	})
	return i
}

// Required checks if the item’s value is non-zero.
func (i *Item) Required(message string) *Item {
	i.Rules = append(i.Rules, &Rule{
		Func: func(value interface{}) (bool, error) {
			switch value := value.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return value != 0, nil
			case string:
				return len(value) > 0, nil
			case time.Time:
				return !value.IsZero(), nil
			case []string:
				return len(value) > 0, nil
			}
			return false, fmt.Errorf("validation.Item.Required: unsupported value type %T", value)
		},
		Args:    []interface{}{RuleTypeRequired},
		Message: message,
		Type:    RuleTypeRequired,
	})
	return i
}

// Validate checks if the item’s value is valid according to the specified
// validation rules. If it is valid, the function returns true. If it is not
// valid, the rule’s validation error message is returned. If an error
// occurred, the error is returned. A returned error does not mean the value is
// invalid, it solely means something went wrong. Rules are checked in order of
// creation. If the item’s value was found to be invalid, any further rules are
// not checked.
func (i *Item) Validate() (bool, string, error) {
	for _, rule := range i.Rules {
		if isValid, err := rule.Func(i.value); err != nil {
			return false, "", err
		} else if !isValid {
			return false, rule.Message, nil
		}
	}
	return true, "", nil
}

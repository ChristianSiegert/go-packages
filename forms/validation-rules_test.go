package forms

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationRule_Validate(t *testing.T) {
	customTestError := errors.New("custom test error")

	fn := func(value string) (bool, error) {
		if value != "foo" {
			return false, nil
		}
		return true, nil
	}

	fn2 := func(value string) (bool, error) {
		return false, customTestError
	}

	tests := []struct {
		rule          *ValidationRule
		value         string
		expectedValid bool
		expectedErr   error
	}{
		{
			rule:          &ValidationRule{Rule: EmailAddress()},
			value:         strings.Repeat("ä", 243) + "@example.com",
			expectedValid: false,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: EmailAddress()},
			value:         "foo@example.com",
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: &Rule{Func: fn}},
			value:         "foo",
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: &Rule{Func: fn}},
			value:         "bar",
			expectedValid: false,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: &Rule{Func: fn2}},
			value:         "baz",
			expectedValid: false,
			expectedErr:   customTestError,
		},
		{
			rule:          &ValidationRule{Rule: MaxLength(10)},
			value:         strings.Repeat("ä", 10),
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: MaxLength(10)},
			value:         strings.Repeat("ä", 11),
			expectedValid: false,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: MinLength(5)},
			value:         strings.Repeat("ä", 4),
			expectedValid: false,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: MinLength(5)},
			value:         strings.Repeat("ä", 5),
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: Required()},
			value:         "",
			expectedValid: false,
			expectedErr:   nil,
		},
		{
			rule:          &ValidationRule{Rule: Required()},
			value:         "hello",
			expectedValid: true,
			expectedErr:   nil,
		},
	}

	for i, test := range tests {
		if result, err := test.rule.Validate(test.value); err != test.expectedErr {
			t.Errorf("Test %d: rule.Validate(%q) returned error %q, expected error %q.", i+1, test.value, err, test.expectedErr)
		} else if result != test.expectedValid {
			t.Errorf("Test %d: rule.Validate(%q) returned %t, expected %t.", i+1, test.value, result, test.expectedValid)
		}
	}
}

func TestEmailAddress(t *testing.T) {
	tests := []struct {
		value         string
		expectedValid bool
		expectedErr   error
	}{
		{strings.Repeat("ä", 243) + "@example.com", false, nil},
		{strings.Repeat("ä", 242) + "@example.com", true, nil},
		{"foo@example.com", true, nil},
		{"@example.com", false, nil},
		{"example.com", false, nil},
		{"a@b", true, nil},
		{"a@", false, nil},
		{"@b", false, nil},
		{"a", false, nil},
		{"@", false, nil},
		{"", false, nil},
	}

	for i, test := range tests {
		if result, err := EmailAddress().Func(test.value); err != test.expectedErr {
			t.Errorf("Test %d: EmailAddress()(%q) returned error %q, expected error %q.", i+1, test.value, err, test.expectedErr)
		} else if result != test.expectedValid {
			t.Errorf("Test %d: EmailAddress()(%q) returned %t, expected %t.", i+1, test.value, result, test.expectedValid)
		}
	}
}

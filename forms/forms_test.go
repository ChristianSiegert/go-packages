package forms

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/ChristianSiegert/go-packages/html/elements"
)

func TestForm_Email(t *testing.T) {
	request1, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}

	request2, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}
	request2.Form = map[string][]string{
		"foo": {"foo@example.com"},
	}

	form1 := New(request1)
	form2 := New(request2)
	form2.ValidationRuleMap = map[string][]*ValidationRule{
		"foo": []*ValidationRule{
			&ValidationRule{
				Rule: Required(),
			},
			&ValidationRule{
				Rule: MaxLength(60),
			},
			&ValidationRule{
				Rule: MinLength(12),
			},
		},
	}

	tests := []struct {
		form        *Form
		name        string
		placeholder string
		expected    *elements.Element
	}{
		{
			form:        form1,
			name:        "foo",
			placeholder: "",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":        "foo",
					"maxlength": "254",
					"name":      "foo",
					"type":      "email",
				},
				Tag: "input",
			},
		},
		{
			form:        form2,
			name:        "foo",
			placeholder: "bar",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":          "foo",
					"maxlength":   "60",
					"minlength":   "12",
					"name":        "foo",
					"placeholder": "bar",
					"required":    "",
					"type":        "email",
					"value":       "foo@example.com",
				},
				Tag: "input",
			},
		},
	}

	for i, test := range tests {
		if result := test.form.Email(test.name, test.placeholder); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %d returned\n%+v\nexpected\n%+v", i+1, result, test.expected)
		}
	}
}

func TestForm_Error(t *testing.T) {
	request, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}
	form1 := New(request)
	form1.ErrorMessageMap = map[string]string{
		"foo": "foo error",
	}

	tests := []struct {
		form      *Form
		fieldName string
		expected  *elements.Element
	}{
		{
			form:      form1,
			fieldName: "foo",
			expected: &elements.Element{
				Attributes: map[string]string{
					"class": "validation-error",
				},
				HasEndTag: true,
				Tag:       "div",
				Text:      "foo error",
			},
		},
		{
			form:      form1,
			fieldName: "bar",
			expected:  nil,
		},
	}

	for i, test := range tests {
		if result := test.form.Error(test.fieldName); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %d form.Error(%q) returned\n%s\nexpected\n%s", i+1, test.fieldName, result, test.expected)
		}
	}
}

func TestForm_HasError(t *testing.T) {
	form := New(nil)
	form.ErrorMessageMap = map[string]string{
		"foo": "foo error",
	}

	tests := []struct {
		form      *Form
		fieldName string
		expected  bool
	}{
		{
			form:      form,
			fieldName: "foo",
			expected:  true,
		},
		{
			form:      form,
			fieldName: "bar",
			expected:  false,
		},
	}

	for _, test := range tests {
		if result := test.form.HasError(test.fieldName); result != test.expected {
			t.Errorf("form.HasError(%q) returned %t, expected %t.", test.fieldName, result, test.expected)
		}
	}
}

func TestForm_Input(t *testing.T) {
	request1, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}

	request2, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}
	request2.Form = map[string][]string{
		"foo": {"Hello, world!"},
	}

	form1 := New(request1)
	form2 := New(request2)
	form2.ValidationRuleMap = map[string][]*ValidationRule{
		"foo": []*ValidationRule{
			&ValidationRule{
				Rule: Required(),
			},
			&ValidationRule{
				Rule: MaxLength(80),
			},
			&ValidationRule{
				Rule: MinLength(3),
			},
		},
	}

	tests := []struct {
		form        *Form
		fieldName   string
		placeholder string
		expected    *elements.Element
	}{
		{
			form:        form1,
			fieldName:   "foo",
			placeholder: "",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":   "foo",
					"name": "foo",
				},
				Tag: "input",
			},
		},
		{
			form:        form2,
			fieldName:   "foo",
			placeholder: "bar",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":          "foo",
					"maxlength":   "80",
					"minlength":   "3",
					"name":        "foo",
					"placeholder": "bar",
					"required":    "",
					"value":       "Hello, world!",
				},
				Tag: "input",
			},
		},
	}

	for i, test := range tests {
		if result := test.form.Input(test.fieldName, test.placeholder); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %d returned\n%+v\nexpected\n%+v", i+1, result, test.expected)
		}
	}
}

func TestForm_Password(t *testing.T) {
	request1, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}

	request2, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}
	request2.Form = map[string][]string{
		"foo": {"password123"},
	}

	form1 := New(request1)
	form2 := New(request2)
	form2.ValidationRuleMap = map[string][]*ValidationRule{
		"foo": []*ValidationRule{
			&ValidationRule{
				Rule: Required(),
			},
			&ValidationRule{
				Rule: MaxLength(50),
			},
			&ValidationRule{
				Rule: MinLength(10),
			},
		},
	}

	tests := []struct {
		form        *Form
		name        string
		placeholder string
		expected    *elements.Element
	}{
		{
			form:        form1,
			name:        "foo",
			placeholder: "",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":   "foo",
					"name": "foo",
					"type": "password",
				},
				Tag: "input",
			},
		},
		{
			form:        form2,
			name:        "foo",
			placeholder: "bar",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":          "foo",
					"maxlength":   "50",
					"minlength":   "10",
					"name":        "foo",
					"placeholder": "bar",
					"required":    "",
					"type":        "password",
					"value":       "password123",
				},
				Tag: "input",
			},
		},
	}

	for i, test := range tests {
		if result := test.form.Password(test.name, test.placeholder); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %d returned\n%+v\nexpected\n%+v", i+1, result, test.expected)
		}
	}
}

func TestForm_Textarea(t *testing.T) {
	request1, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}

	request2, err := http.NewRequest("GET", "/", &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Creating request failed unexpectedly: %s", err)
	}
	request2.Form = map[string][]string{
		"foo": {"Hello, world!"},
	}

	form1 := New(request1)
	form2 := New(request2)
	form2.ValidationRuleMap = map[string][]*ValidationRule{
		"foo": []*ValidationRule{
			&ValidationRule{
				Rule: Required(),
			},
			&ValidationRule{
				Rule: MaxLength(1000),
			},
			&ValidationRule{
				Rule: MinLength(100),
			},
		},
	}

	tests := []struct {
		form        *Form
		name        string
		placeholder string
		expected    *elements.Element
	}{
		{
			form:        form1,
			name:        "foo",
			placeholder: "",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":   "foo",
					"name": "foo",
				},
				HasEndTag: true,
				Tag:       "textarea",
			},
		},
		{
			form:        form2,
			name:        "foo",
			placeholder: "bar",
			expected: &elements.Element{
				Attributes: map[string]string{
					"id":          "foo",
					"maxlength":   "1000",
					"minlength":   "100",
					"name":        "foo",
					"placeholder": "bar",
					"required":    "",
				},
				HasEndTag: true,
				Tag:       "textarea",
				Text:      "Hello, world!",
			},
		},
	}

	for i, test := range tests {
		if result := test.form.Textarea(test.name, test.placeholder); !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Test %d returned\n%+v\nexpected\n%+v", i+1, result, test.expected)
		}
	}
}

func TestForm_Validate(t *testing.T) {
	customTestError := errors.New("test: foo already taken")

	tests := []struct {
		requestForm             map[string][]string
		validationRuleMap       map[string][]*ValidationRule
		expectedValid           bool
		expectedErr             error
		expectedErrorMessageMap map[string]string
	}{
		{
			requestForm:             nil,
			validationRuleMap:       nil,
			expectedValid:           false,
			expectedErr:             ErrNoValidationRules,
			expectedErrorMessageMap: nil,
		},
		{
			requestForm: map[string][]string{
				"foo": {strings.Repeat("ä", 5)},
			},
			validationRuleMap: map[string][]*ValidationRule{
				"foo": []*ValidationRule{
					&ValidationRule{
						Message: "Foo is already taken.",
						Rule: &Rule{
							Func: func(value string) (bool, error) {
								return false, customTestError
							},
						},
					},
				},
			},
			expectedValid:           false,
			expectedErr:             customTestError,
			expectedErrorMessageMap: nil,
		},
		{
			requestForm: map[string][]string{
				"foo": {strings.Repeat("ä", 5)},
				"bar": {strings.Repeat("ä", 10)},
			},
			validationRuleMap: map[string][]*ValidationRule{
				"foo": []*ValidationRule{
					&ValidationRule{
						Message: "Foo must have at least 5 characters.",
						Rule:    MinLength(5),
					},
					&ValidationRule{
						Message: "Foo must have at most 50 characters",
						Rule:    MaxLength(50),
					},
				},
				"bar": []*ValidationRule{
					&ValidationRule{
						Message: "Bar must have at least 5 characters.",
						Rule:    MinLength(5),
					},
					&ValidationRule{
						Message: "Bar must have at most 50 characters",
						Rule:    MaxLength(50),
					},
				},
			},
			expectedValid:           true,
			expectedErr:             nil,
			expectedErrorMessageMap: nil,
		},
		{
			requestForm: map[string][]string{
				"foo": {strings.Repeat("ä", 4)},
				"bar": {strings.Repeat("ä", 10)},
			},
			validationRuleMap: map[string][]*ValidationRule{
				"foo": []*ValidationRule{
					&ValidationRule{
						Message: "Foo is required.",
						Rule:    Required(),
					},
					&ValidationRule{
						Message: "Foo must have at least 5 characters.",
						Rule:    MinLength(5),
					},
					&ValidationRule{
						Message: "Foo must have at most 50 characters",
						Rule:    MaxLength(50),
					},
				},
				"bar": []*ValidationRule{
					&ValidationRule{
						Message: "Bar is required.",
						Rule:    Required(),
					},
					&ValidationRule{
						Message: "Bar must have at least 5 characters.",
						Rule:    MinLength(5),
					},
					&ValidationRule{
						Message: "Bar must have at most 50 characters",
						Rule:    MaxLength(50),
					},
				},
				"baz": []*ValidationRule{
					&ValidationRule{
						Message: "Baz is required.",
						Rule:    Required(),
					},
				},
			},
			expectedValid: false,
			expectedErr:   nil,
			expectedErrorMessageMap: map[string]string{
				"foo": "Foo must have at least 5 characters.",
				"baz": "Baz is required.",
			},
		},
	}

	for i, test := range tests {
		request, err := http.NewRequest("GET", "/", &bytes.Buffer{})
		if err != nil {
			t.Fatalf("Creating request failed unexpectedly: %s", err)
		}

		request.Form = test.requestForm
		form := New(request)
		form.ValidationRuleMap = test.validationRuleMap

		expectedForm := New(request)
		expectedForm.ValidationRuleMap = test.validationRuleMap
		expectedForm.ErrorMessageMap = test.expectedErrorMessageMap

		if result, err := form.Validate(); err != test.expectedErr {
			t.Errorf("Test %d: form.Validate() returned error %q, expected error %q.", i+1, err, test.expectedErr)
		} else if !reflect.DeepEqual(form, expectedForm) {
			t.Errorf("Test %d: form and expected form differ:\n%+v\n%+v", i+1, form, expectedForm)
		} else if !reflect.DeepEqual(result, test.expectedValid) {
			t.Errorf("Test %d form.Validate() returned %t, expected %t.", i+1, result, test.expectedValid)
		}
	}
}

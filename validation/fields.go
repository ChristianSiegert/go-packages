// Package validation provides means to validate values.
package validation

// Messages is a map whose keys are field names and whose values are validation
// error messages. The map only contains the names of fields that failed
// validation.
type Messages map[string]string

// Fields is a map whose keys are field names and whose values are *Field.
type Fields map[string]*Field

// Add adds a field whose value is to be validated. Validation rules
// must be attached to the field itself.
func (f Fields) Add(fieldName string, value interface{}) *Field {
	field := &Field{
		value: value,
	}

	f[fieldName] = field
	return field
}

// Validate validates all fields.
func (f Fields) Validate() (Messages, error) {
	var messages Messages

	for fieldName, field := range f {
		if isValid, message, err := field.Validate(); err != nil {
			return nil, err
		} else if !isValid {
			if messages == nil {
				messages = make(Messages)
			}
			messages[fieldName] = message
		}
	}

	return messages, nil
}

// Package validation provides validation for values.
package validation

// Items manages Item objects.
type Items map[string]*Item

// New returns a new instance of Items.
func New() Items {
	return make(Items)
}

// Add adds an item whose value is to be validated. Validation rules must be
// attached to the item itself.
func (i Items) Add(name string, value interface{}) *Item {
	item := &Item{
		value: value,
	}

	i[name] = item
	return item
}

// Validate validates all items.
func (i Items) Validate() (Messages, error) {
	var messages Messages

	for name, item := range i {
		if isValid, message, err := item.Validate(); err != nil {
			return nil, err
		} else if !isValid {
			if messages == nil {
				messages = make(Messages)
			}
			messages[name] = message
		}
	}

	return messages, nil
}

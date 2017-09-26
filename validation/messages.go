package validation

import "fmt"

// Messages is a map whose keys are item names and whose values are validation
// error messages. The map only contains the names of items that failed
// validation.
type Messages map[string]string

// Error implements the Error interface.
func (m Messages) Error() string {
	return fmt.Sprintf("%#v", m)
}

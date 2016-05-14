package validation

// Rule contains the validation function and information about it.
type Rule struct {
	// Arguments that Func was called with.
	Args []interface{}

	// Func returns whether the argument is valid, or that an error occurred
	// while validating. A returned error does not mean the argument is invalid,
	// it solely means something went wrong while validating.
	Func func(interface{}) (bool, error)

	// Message that informs the user if her input is invalid.
	Message string

	// Type gives information about the rule type, e.g. RuleTypeMaxLength means
	// it is a rule for checking maximum length. A value of 0 means no type is
	// provided.
	Type int
}

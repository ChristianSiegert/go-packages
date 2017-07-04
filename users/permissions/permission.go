package permissions

import (
	"encoding"
)

// Permission is an action that a user is allowed to perform.
type Permission interface {
	// Name returns the permission’s name.
	Name() string

	encoding.TextMarshaler
}

type permission string

// NewPermission returns a new instance of an unexported type that implements
// the Permission interface. It is a convenience function to avoid having to
// create a type whenever the Permission interface must be satisfied.
func NewPermission(name string) Permission {
	return permission(name)
}

func (p permission) MarshalText() ([]byte, error) {
	return []byte(p.Name()), nil
}

func (p permission) Name() string {
	return string(p)
}

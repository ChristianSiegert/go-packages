package permissions

// Permission that can be added to a role.
type Permission interface {
	Name() string
}

type permission string

// NewPermission returns a new instance of an unexported type that implements
// the Permission interface. It is a convenience function to avoid having to
// create a type whenever the Permission interface must be satisfied.
func NewPermission(name string) Permission {
	return permission(name)
}

func (p permission) Name() string {
	return string(p)
}

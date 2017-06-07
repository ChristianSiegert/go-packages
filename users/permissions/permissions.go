package permissions

type Permission interface {
	Name() string
}

type permission string

// New returns a new instance of an unexported type that implements the
// Permission interface. It is a convenience function to avoid having to create
// a type whenever the Permission interface must be satisfied.
func New(name string) Permission {
	return permission(name)
}

func (p permission) Name() string {
	return string(p)
}

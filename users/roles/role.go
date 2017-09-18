package roles

// Role of a user.
type Role struct {
	ID          int
	Name        string
	Permissions Permissions
}

// New returns a new Role.
func New(id int, name string, permissions ...Permission) *Role {
	perms := make(Permissions, len(permissions))
	perms.Add(permissions...)

	return &Role{
		ID:          id,
		Name:        name,
		Permissions: perms,
	}
}

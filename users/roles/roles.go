// Package roles provides an interface and implementation for user roles.
package roles

import "github.com/ChristianSiegert/go-packages/users/permissions"

// Role of a user.
type Role interface {
	// Name returns the role’s name.
	Name() string

	// Permissions returns all permissions the role has been granted.
	Permissions() permissions.Map

	// SetName sets the role’s name.
	SetName(name string)
}

// role is an unexported implementation of the Role interface.
type role struct {
	name        string
	permissions permissions.Map
}

// New returns a new instance of an unexported type that implements the
// Role interface.
func New(name string, permissions permissions.Map) Role {
	return &role{
		name:        name,
		permissions: permissions,
	}
}

// Name returns the role’s name.
func (r *role) Name() string {
	return r.name
}

// Permissions returns all permissions the role has been granted.
func (r *role) Permissions() permissions.Map {
	return r.permissions
}

func (r *role) SetName(name string) {
	r.name = name
}

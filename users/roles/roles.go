// Package roles provides an interface and implementation for user roles.
package roles

import "github.com/ChristianSiegert/go-packages/users/permissions"

// Role a user can have.
type Role interface {
	// HasPermission returns whether the role has the given permission.
	HasPermission(permissions.Permission) bool

	// Name returns the role’s name.
	Name() string

	// Permissions returns all permissions the role has been granted.
	Permissions() []permissions.Permission
}

// role is an unexported implementation of the Role interface.
type role struct {
	name        string
	permissions []permissions.Permission
}

// New returns a new instance of an unexported type that implements the
// Role interface.
func New(name string, permissions ...permissions.Permission) Role {
	return &role{
		name:        name,
		permissions: permissions,
	}
}

// HasPermission returns whether the role has the given permission.
func (r role) HasPermission(permission permissions.Permission) bool {
	for _, perm := range r.permissions {
		if permission == perm {
			return true
		}
	}
	return false
}

// Name returns the role’s name.
func (r role) Name() string {
	return r.name
}

// Permissions returns all permissions the role has been granted.
func (r role) Permissions() []permissions.Permission {
	return r.permissions
}

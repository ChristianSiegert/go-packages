// Package roles provides an interface and implementation for user roles.
package roles

import "github.com/ChristianSiegert/go-packages/users/permissions"
import "encoding/json"

// Role of a user.
type Role interface {
	ID() int

	// Name returns the role’s name.
	Name() string

	// Permissions returns all permissions the role has been granted.
	Permissions() permissions.Map

	// SetName sets the role’s name.
	SetName(name string)
}

// role is an unexported implementation of the Role interface.
type role struct {
	id          int
	name        string
	permissions permissions.Map
}

// jsonRole is an unexported type that is used to JSON encode and decode a role.
type jsonRole struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Permissions permissions.Map `json:"permissions"`
}

// New returns a new instance of an unexported type that implements the
// Role interface.
func New(id int, name string, permissions permissions.Map) Role {
	return &role{
		id:          id,
		name:        name,
		permissions: permissions,
	}
}

// ID returns the role’s ID.
func (r *role) ID() int {
	return r.id
}

func (r *role) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonRole{
		ID:          r.ID(),
		Name:        r.Name(),
		Permissions: r.Permissions(),
	})
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

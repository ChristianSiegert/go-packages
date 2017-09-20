package roles

import (
	"encoding/json"
)

// Permissions contains permissions that have been granted.
type Permissions map[Permission]bool

// Add adds permissions.
func (p Permissions) Add(permissions ...Permission) {
	for _, permission := range permissions {
		p[permission] = true
	}
}

// Has returns whether permission exists in the map.
func (p Permissions) Has(permission Permission) bool {
	_, isPresent := p[permission]
	return isPresent
}

// HasOne returns whether one of the provided permissions exists in the map.
func (p Permissions) HasOne(permissions ...Permission) bool {
	for _, permission := range permissions {
		if _, isPresent := p[permission]; isPresent {
			return true
		}
	}
	return false
}

// MarshalJSON JSON encodes the map.
func (p Permissions) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[Permission]bool(p))
}

// Remove removes permissions.
func (p Permissions) Remove(permissions ...Permission) {
	for _, permission := range permissions {
		delete(p, permission)
	}
}

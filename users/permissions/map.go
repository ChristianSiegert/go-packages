package permissions

// Map contains permissions that have been granted.
type Map map[Permission]bool

// NewMap returns a new instance of Map.
func NewMap(permissions ...Permission) Map {
	m := make(Map, len(permissions))
	m.Add(permissions...)
	return m
}

// Add adds permissions.
func (m Map) Add(permissions ...Permission) {
	for _, permission := range permissions {
		m[permission] = true
	}
}

// Has returns whether a permission has been added.
func (m Map) Has(permission Permission) bool {
	_, present := m[permission]
	return present
}

// Remove removes permissions.
func (m Map) Remove(permissions ...Permission) {
	for _, permission := range permissions {
		delete(m, permission)
	}
}

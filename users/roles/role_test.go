package roles

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	permissionPostsCreate := Permission("posts.create")
	permissionUsersCreate := Permission("users.create")
	permissionUsersDelete := Permission("users.delete")

	tests := []struct {
		id          int
		name        string
		permissions []Permission
	}{
		{
			id:   1,
			name: "administrator",
			permissions: []Permission{
				permissionPostsCreate,
				permissionUsersCreate,
				permissionUsersDelete,
			},
		},
		{
			id:   2,
			name: "user",
			permissions: []Permission{
				permissionPostsCreate,
			},
		},
	}

	for _, test := range tests {
		role := New(test.id, test.name, test.permissions...)

		if role.Name != test.name {
			t.Errorf("Expected name %q, got %q.", test.name, role.Name)
		} else if reflect.DeepEqual(role.Permissions, test.permissions) {
			t.Errorf("Expected permissions %v, got %v.", test.permissions, role.Permissions)
		}
	}
}

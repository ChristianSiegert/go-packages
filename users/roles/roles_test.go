package roles

import (
	"reflect"
	"testing"

	"github.com/ChristianSiegert/go-packages/users/permissions"
)

func TestNew(t *testing.T) {
	permissionPostsCreate := permissions.NewPermission("posts.create")
	permissionUsersCreate := permissions.NewPermission("users.create")
	permissionUsersDelete := permissions.NewPermission("users.delete")

	tests := []struct {
		name        string
		permissions permissions.Map
	}{
		{
			name: "administrator",
			permissions: permissions.NewMap(
				permissionPostsCreate,
				permissionUsersCreate,
				permissionUsersDelete,
			),
		},
		{
			name: "user",
			permissions: permissions.NewMap(
				permissionPostsCreate,
			),
		},
	}

	for _, test := range tests {
		role := New(test.name, test.permissions)

		if role.Name() != test.name {
			t.Errorf("Expected name %q, got %q.", test.name, role.Name())
		} else if reflect.DeepEqual(role.Permissions, test.permissions) {
			t.Errorf("Expected permissions %v, got %v.", test.permissions, role.Permissions())
		}
	}
}

func TestSetName(t *testing.T) {
	test := struct {
		name         string
		expectedName string
	}{
		name:         "bar",
		expectedName: "bar",
	}

	role := New("foo", nil)
	role.SetName(test.name)

	if role.Name() != test.expectedName {
		t.Errorf("Expected name %q, got %q.", test.expectedName, role.Name())
	}
}

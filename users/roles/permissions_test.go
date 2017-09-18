package roles

import (
	"testing"
)

var permissionA = Permission("A")
var permissionB = Permission("B")

func TestPermissions_Has(t *testing.T) {
	permissions := make(Permissions)
	permissions.Add(permissionA)
	permissions.Add(permissionB)

	if len(permissions) != 2 || !permissions.Has(permissionA) || !permissions.Has(permissionB) {
		t.Error("Has failed.")
	}
}

func TestPermissions_Remove(t *testing.T) {
	permissions := make(Permissions)
	permissions.Add(permissionA)
	permissions.Add(permissionB)
	permissions.Remove(permissionA)

	if len(permissions) != 1 || !permissions.Has(permissionB) {
		t.Error("Removing permission failed.")
	}
}

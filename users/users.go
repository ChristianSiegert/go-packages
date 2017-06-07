package users

import (
	"github.com/ChristianSiegert/go-packages/users/roles"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user.
type User interface {
	EmailAddress() string
	Id() string
	Name() string
	PasswordHash() []byte
	Role() roles.Role
	Username() string
}

// HashPassword hashes a password with bcrypt. cost is in interval
// [bcrypt.MinCost, bcrypt.MaxCost], i.e. [4, 31]. Use bcrypt.DefaultCost, or
// 10, if unsure.
func HashPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// IsPassword returns whether password is the plaintext equivalent of hash.
func IsPassword(password, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

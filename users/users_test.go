package users

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password           []byte
		cost               int
		expectedHashLength int
		expectedError      error
	}{
		{
			password:           []byte("mypass"),
			cost:               bcrypt.MinCost - 1,
			expectedHashLength: 60,
			expectedError:      nil,
		},
		{
			password:           []byte("mypass"),
			cost:               bcrypt.MinCost,
			expectedHashLength: 60,
			expectedError:      nil,
		},
		{
			password:           []byte("mypass"),
			cost:               bcrypt.MaxCost + 1,
			expectedHashLength: 0,
			expectedError:      bcrypt.InvalidCostError(bcrypt.MaxCost + 1),
		},
	}

	for _, test := range tests {
		if hash, err := HashPassword(test.password, test.cost); err != test.expectedError {
			t.Errorf("Expected error %q, got %q.", test.expectedError, err)
		} else if length := len(hash); length != test.expectedHashLength {
			t.Errorf("Expected hash length %d, got %d.", test.expectedHashLength, length)
		}
	}
}

func TestIsPassword(t *testing.T) {
	var tests = []struct {
		password       []byte
		hash           []byte
		expectedResult bool
		expectedError  error
	}{
		{
			password:       []byte("foo"),
			hash:           []byte("abc"),
			expectedResult: false,
			expectedError:  bcrypt.ErrHashTooShort,
		},
		{
			password:       []byte("foo"),
			hash:           []byte("$11111222222333333444444555555666666777777888888999999000000"),
			expectedResult: false,
			expectedError:  nil,
		},
		{
			password:       []byte("foo"),
			hash:           []byte("$2a$10$lfIzGe3pc.4ip53cChrrouusqRGBpj523la/jNKWalKmS80f3AXbW"),
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		if result, err := IsPassword(test.password, test.hash); err != test.expectedError {
			t.Errorf("Expected error %q, got %q.", test.expectedError, err)
		} else if result != test.expectedResult {
			t.Errorf("Expected result %t, got %t.", test.expectedResult, result)
		}
	}
}

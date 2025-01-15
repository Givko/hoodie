package security_test

import (
	"testing"

	"github.com/plamendelchev/hoodie/internal/security"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashedPassword, err := security.HashPassword(password)

	if err != nil {
		t.Errorf("An err ocucred while hashing the password: %v", err)
	}

	if len(hashedPassword.Password) == 0 {
		t.Errorf("The hashed password is empty")
	}

	if len(hashedPassword.Salt) == 0 {
		t.Errorf("The salt is empty")
	}

	areEqual := security.ComparePassword(password, hashedPassword.Password, hashedPassword.Salt)
	if !areEqual {
		t.Errorf("The passwords do not match")
	}
}

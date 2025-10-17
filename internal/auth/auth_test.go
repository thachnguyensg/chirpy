package auth

import (
	"testing"
)

func TestHashPasswordOK(t *testing.T) {
	password := "my_secure_password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if !match {
		t.Fatalf("Password does not match hash")
	}
}

func TestHashPasswordWrong(t *testing.T) {
	password := "my_secure_password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	match, err := CheckPasswordHash("wrong_password", hash)
	if err != nil {
		t.Fatalf("Error checking password hash: %v", err)
	}
	if match {
		t.Fatalf("Password should not match hash")
	}
}

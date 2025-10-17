package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func TestMakeAndValidateJWTOK(t *testing.T) {
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Error generating UUID: %v", err)
	}
	secret := "my_secret_key"
	token, err := MakeJWT(userID, secret, time.Minute*1)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}
	if returnedUserID != userID {
		t.Fatalf("Returned user ID does not match original. Got %s, want %s", returnedUserID, userID)
	}
}

func TestMakeAndValidateJWTWrongKey(t *testing.T) {
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Error generating UUID: %v", err)
	}
	secret := "my_secret_key"
	token, err := MakeJWT(userID, secret, time.Minute*1)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, "wrong_secret_key")
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Fatalf("Expected error '%v' for wrong key, got: %v", jwt.ErrTokenSignatureInvalid, err)
	}
}

func TestMakeAndValidateJWTExpired(t *testing.T) {
	userID, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Error generating UUID: %v", err)
	}
	secret := "my_secret_key"
	token, err := MakeJWT(userID, secret, -time.Minute*1)
	if err != nil {
		t.Fatalf("Error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("Expected error '%v' for expired token, got: %v", jwt.ErrTokenExpired, err)
	}
}

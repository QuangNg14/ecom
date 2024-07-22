package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be not empty")
	}

	if hash == password {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("expected password and hash to match")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("expected password and hash to not match")
	}
}

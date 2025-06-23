package auth

import (
	"testing"
)

func TestPasswordMatching(t *testing.T) {
	password := "MIamigo"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Error matching password to hash")
	}
}

func TestPasswordMatchingFalse(t *testing.T) {
	hash, err := HashPassword("Miamigo")
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash("MIamigo", hash)
	if err == nil {
		t.Errorf("No error with non matching password/hash")
	}
}

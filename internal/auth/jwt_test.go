package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
)

func TestJwtToken(t *testing.T) {
	tokenSecret := "aladsoo234adlafk"
	userID := uuid.New()
	expiresIn := 10 * time.Second 

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("couldnt make JWT: %v", err)
	}

	id, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("Couldnt validate JWT: %v", err)
	}

	if id != userID {
		t.Errorf("Non matching id(%v) and userID(%v)", id, userID)
	}
}

func TestJwtToken1(t *testing.T) {
	tokenSecret := "aladsoo234adlafk"
	userID := uuid.New()
	expiresIn := time.Duration(0) 

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("couldnt make JWT: %v", err)
	}

	id, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Logf("Couldnt validate JWT: %v", err)
		return
	}

	if id == userID {
		t.Errorf("expired token worked")
	}
}

func TestJwtToken2(t *testing.T) {
	tokenSecret := "aladsoo234adlafk"
	userID := uuid.New()
	expiresIn := 10 * time.Second 

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("couldnt make JWT: %v", err)
	}

	id, err := ValidateJWT(tokenString, "aladsoo")
	if err != nil {
		t.Logf("Couldnt validate JWT: %v", err)
		return
	}

	if id == userID {
		t.Errorf("Validated JWT with wrong tokenSecret")
	}
}





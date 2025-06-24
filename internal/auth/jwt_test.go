package auth

import (
	"testing"
	"time"
	"net/http"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
    tests := []struct {
        name      string
        headers   http.Header
        wantToken string
        wantErr   bool
    }{
        {
            name: "Valid Bearer token",
            headers: http.Header{
                "Authorization": []string{"Bearer abcdef12345"},
            },
            wantToken: "abcdef12345",
            wantErr:   false,
        },
        {
            name:    "Missing Authorization header",
            headers: http.Header{},
            wantToken: "",
            wantErr:   true,
        },
        {
            name: "Authorization header without Bearer",
            headers: http.Header{
                "Authorization": []string{"abcdef12345"},
            },
            wantToken: "abcdef12345", 
            wantErr:   false,
        },
        {
            name: "Bearer token with extra spaces",
            headers: http.Header{
                "Authorization": []string{"Bearer    tokenwithspaces"},
            },
            wantToken: "   tokenwithspaces",
            wantErr:   false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotToken, err := GetBearerToken(tt.headers)

            if (err != nil) != tt.wantErr {
                t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
            }

            if gotToken != tt.wantToken {
                t.Errorf("expected token: %q, got: %q", tt.wantToken, gotToken)
            }
        })
    }
}

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





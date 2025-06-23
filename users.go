package main

import (
	"encoding/json"
	"net/http"
	"time"
	//"regexp"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	/*
	if !isEmail(params.Body) {
		respondWithError(w, 400, "invalid email", nil)
		return
	}
	*/

	user, err := cfg.queries.CreateUser(req.Context(), params.Email)
	if err != nil {
		s := fmt.Sprintf("Error creating user (email: %s)", params.Email)
		respondWithError(w, 400, s, err)
		return
	}

	respondWithJson(w, 201, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
	
}

/*
func isEmail(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
  return emailRegex.MatchString(s)
}
*/

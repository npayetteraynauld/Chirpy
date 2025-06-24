package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/npayetteraynauld/Chirpy/internal/database"
	"github.com/npayetteraynauld/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, req *http.Request) {
	//retrieve tokenstring from header
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authentication header", err)
		return
	}
	
	//Validate JWT
	userIDFromToken, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid auth token", err)
		return 
	}


	type parameters struct {
		Body   string    `json:"body"`
	}

	type returnVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}


	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.queries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: userIDFromToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}
	respondWithJson(w, 201, returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: cleanString(chirp.Body),
		UserID: chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.queries.GetChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	type returnVals struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
	}

	var returns []returnVals
	for _, chirp := range chirps {
		returns = append(returns, returnVals{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, returns)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	type returnVals struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
	}

	chirpId := req.PathValue("chirpID")

	id, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse id", err)
	}

	chirp, err := cfg.queries.GetChirp(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJson(w, http.StatusOK, returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func cleanString(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitString := strings.Split(s, " ")

	for i, word := range splitString {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				splitString[i] = "****"
				break
			}
		}
	}

	return strings.Join(splitString, " ")
}

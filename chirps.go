package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"sort"

	"github.com/google/uuid"
	"github.com/npayetteraynauld/Chirpy/internal/database"
	"github.com/npayetteraynauld/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		type returnVals struct {
				ID        uuid.UUID `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Body      string    `json:"body"`
				UserID    uuid.UUID `json:"user_id"`
		}

		var chirps []database.Chirp
		var err error

		//Check for optional queries
		optionalAuthorID := req.URL.Query().Get("author_id")
		optionalSorting := req.URL.Query().Get("sort")

		if optionalAuthorID != "" {
			//Get only chirps from author
			authorID, err := uuid.Parse(optionalAuthorID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't parse ID", err)
				return
			}

			chirps, err = cfg.queries.GetChirpsFromUserID(req.Context(), authorID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
				return 
			}

		} else {
			//Get all chirps
			chirps, err = cfg.queries.GetChirps(req.Context())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
				return
			}
		}

		//Check if optional sort arg is "desc"
		if optionalSorting == "desc" {
			sort.Slice(chirps, func(i, j int) bool{
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
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

	} else if req.Method == "POST" {
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

	}  else {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
	}
}

func (cfg *apiConfig) handlerChirpsByID(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
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

	} else if req.Method == "DELETE" {
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

		chirpID := req.PathValue("chirpID")
		
		id, err := uuid.Parse(chirpID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse id", err)
			return
		}
		
		chirp, err := cfg.queries.GetChirp(req.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
			return
		}

		if chirp.UserID != userIDFromToken {
			respondWithError(w, http.StatusForbidden, "Not permitted to Delete chirp", nil)
			return
		}

		err = cfg.queries.DeleteChirp(req.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Couldn't delete chirp", err)
			return
		}

		w.WriteHeader(204)

	} else {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
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


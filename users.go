package main

import (
	"encoding/json"
	"net/http"
	"time"
	//"regexp"
	"fmt"

	"github.com/google/uuid"
	"github.com/npayetteraynauld/Chirpy/internal/auth"
	"github.com/npayetteraynauld/Chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	//check if password was inputted
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid Password", nil)
	}

	//hash password
	hashedpsswd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	//Create user in DB
	user, err := cfg.queries.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedpsswd,
	})
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
		IsChirpyRed: user.IsChirpyRed,
	})
	
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return 
	}

	err = cfg.queries.UpdateEmailAndPassword(req.Context(), database.UpdateEmailAndPasswordParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userIDFromToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	user, err := cfg.queries.GetUserFromID(req.Context(), userIDFromToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch user form db", err)
		return 
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,         
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt, 
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.queries.GetUserFromEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	} 

	expires := time.Hour
	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, expires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return 
	}
	
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make RefreshToken", err)
		return  
	}

	refreshTokenSQL, err := cfg.queries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshTokenString,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't register refresh tokken", err)
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)

	respondWithJson(w, http.StatusOK, User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: tokenString,
			RefreshToken: refreshTokenSQL.Token,
			IsChirpyRed: user.IsChirpyRed,
		})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authentication header", err)
		return
	}

	refreshToken, err := cfg.queries.GetRefreshTokenFromId(req.Context(), authToken)
	if err != nil || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	
	user, err := cfg.queries.GetUserFromRefreshToken(req.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get userID from db", err)
		return 
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't make JWT", err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	type ReturnVals struct {
		Token string `json:"token"`
	}

	respondWithJson(w, http.StatusOK, ReturnVals{Token: token})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No authentication header", err)
		return
	}

	err = cfg.queries.RevokeRefreshToken(req.Context(), authToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error revoking token", err)
		return
	}

	w.WriteHeader(204)
}



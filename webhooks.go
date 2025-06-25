package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/npayetteraynauld/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get apikey", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized request", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return 
	}

	uuid, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt parse user id", err)
		return 
	}

	_, err = cfg.queries.UpgradeUserRed(req.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade specified user", err)
		return
	}

	w.WriteHeader(204)
}

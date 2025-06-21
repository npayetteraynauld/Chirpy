package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	

	respondWithJson(w, http.StatusOK, returnVals{
		CleanedBody: cleanString(params.Body),
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

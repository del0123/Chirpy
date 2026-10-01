package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ChirpRequest struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Instead of Marshal, I could also have used the json.NewEncoder(w).Encode(payload) method (streamed).
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, ErrorResponse{Error: msg})
}

func chirpProfanityCleaner(chirp string) string {
	// TODO: Implement profanity cleaner logic
	words := strings.Split(chirp, " ")

	for i, word := range words {
		lowerWord := strings.ToLower(word)

		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}

	cleanedChirp := chirpProfanityCleaner(body)
	return cleanedChirp, nil
}

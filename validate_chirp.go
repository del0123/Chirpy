package main

import (
	"encoding/json"
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := ChirpRequest{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	params.Body = chirpProfanityCleaner(params.Body)

	respondWithJSON(w, http.StatusOK, ChirpResponse{CleanedBody: params.Body})

}

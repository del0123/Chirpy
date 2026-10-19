package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/del0123/Chirpy/internal/auth"
	"github.com/del0123/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var params struct {
		Body string `json:"body"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenUserID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	validatedBody, err := validateChirp(params.Body)
	if err != nil {
		http.Error(w, "Invalid chirp body", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   validatedBody,
		UserID: tokenUserID,
	})
	if err != nil {
		http.Error(w, "Failed to create chirp", http.StatusInternalServerError)
		return
	}

	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    tokenUserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	inputAuthorID := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error

	if inputAuthorID != "" {
		authorID, err := uuid.Parse(inputAuthorID)
		if err != nil {
			http.Error(w, "Invalid author_id format", http.StatusBadRequest)
			return
		}
		dbChirps, err = cfg.db.GetChirpsFromUser(ctx, authorID)
		if err != nil {
			http.Error(w, "Failed to get chirps from author", http.StatusInternalServerError)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetChirps(ctx)
		if err != nil {
			http.Error(w, "Failed to get chirps", http.StatusInternalServerError)
			return
		}
	}

	chirpsResponse := []Chirp{}

	for _, chirp := range dbChirps {
		chirpResponse := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpsResponse = append(chirpsResponse, chirpResponse)
	}

	respondWithJSON(w, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid Chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirp(ctx, chirpID)
	if err != nil {
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	idString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirp(ctx, chirpID)
	if err != nil {
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	if chirp.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = cfg.db.DeleteChirp(ctx, chirpID)
	if err != nil {
		http.Error(w, "Failed to delete chirp", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

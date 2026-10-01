package main

import (
	"encoding/json"
	"net/http"
	"time"

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

	var params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validatedBody, err := validateChirp(params.Body)
	if err != nil {
		http.Error(w, "Invalid chirp body", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   validatedBody,
		UserID: params.UserID,
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
		UserID:    chirp.UserID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chirps, err := cfg.db.GetChirps(ctx)
	if err != nil {
		http.Error(w, "Failed to get chirps", http.StatusInternalServerError)
		return
	}

	chirpResponses := []Chirp{}

	for _, chirp := range chirps {
		chirpResponse := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpResponses = append(chirpResponses, chirpResponse)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponses)
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

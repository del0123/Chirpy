package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/del0123/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	polkaHeaderKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if polkaHeaderKey != cfg.polkaKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	var request polkaWebhookRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeChirpyRed(ctx, request.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

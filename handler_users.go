package main

import (
	"encoding/json"
	"log"

	"net/http"
	"time"

	"github.com/del0123/Chirpy/internal/auth"
	"github.com/del0123/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(ctx, database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid Username or Password", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.GetUserByEmail(ctx, params.Email)
	if err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	if !match {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error Making JWT token string: %v", err)
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error Making Refresh token string: %v", err)
		http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
		return
	}

	_, err = cfg.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		UserID:    user.ID,
	})
	if err != nil {
		log.Printf("Error Creating refresh token in database: %v", err)
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshTokenString,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(ctx, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwtString, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtString,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.db.RevokeRefreshToken(ctx, tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

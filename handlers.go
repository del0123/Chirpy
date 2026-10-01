package main

import (
	"fmt"
	"net/http"
)

// this function is a health checker
func endpointHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Endpoint is OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	count := cfg.fileserverHits.Load()

	htmlResponse := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", count)

	fmt.Fprint(w, htmlResponse)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if cfg.platform != "dev" {
		http.Error(w, "Unauthorized Access attempt, FBI notified!", http.StatusForbidden)
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "All users deleted and hits reset to %d", cfg.fileserverHits.Load())

}

package main

import (
	"sync/atomic"

	"github.com/del0123/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

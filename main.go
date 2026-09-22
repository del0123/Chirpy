package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := &apiConfig{}

	fileDir := http.Dir(".")
	fileHandler := http.FileServer(fileDir)
	appHandler := http.StripPrefix("/app", fileHandler)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))
	mux.HandleFunc("GET /api/healthz", endpointHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	log.Printf("Serving on port: %s\n", server.Addr)

}

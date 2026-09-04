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

	fileDir := http.Dir(".")
	fileHandler := http.FileServer(fileDir)
	mux.Handle("/", fileHandler)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	log.Printf("Serving on port: %s\n", server.Addr)
}

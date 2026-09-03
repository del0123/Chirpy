package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

	fmt.Println("Server stopped.")
}

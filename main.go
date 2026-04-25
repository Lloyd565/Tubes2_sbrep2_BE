package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"tubes2/backend/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/traverse", api.TraverseHandler)
	mux.HandleFunc("/api/log", api.LogHandler)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)
	fmt.Println("BE running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
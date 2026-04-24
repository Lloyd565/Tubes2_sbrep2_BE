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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "index.html")
	})

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)

	fmt.Println("Backend running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
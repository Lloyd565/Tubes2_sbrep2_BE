package main

import (
	"fmt"
	"net/http"
	"tubes2/backend/internal/api"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/traverse", api.TraverseHandler)

	fmt.Println("Backend running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

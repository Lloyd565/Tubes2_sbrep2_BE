package api

import (
	"encoding/json"
	"net/http"
)

type TraverseRequest struct {
	URL       string `json:"url"`
	HTML      string `json:"html"`
	Algorithm string `json:"algorithm"`
	Selector  string `json:"selector"`
	Limit     int    `json:"limit"`
}

type TraverseResponse struct {
	Tree         interface{} `json:"tree"`
	Matches      interface{} `json:"matches"`
	TraversalLog []string    `json:"traversal_log"`
	Stats        interface{} `json:"stats"`
}

func TraverseHandler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		return
	}

	var req TraverseRequest
	json.NewDecoder(r.Body).Decode(&req)

	// TODO: panggil scraper → parser → traversal

	resp := TraverseResponse{
		TraversalLog: []string{"test ok"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

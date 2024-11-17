package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TestHandler returns a simple string message
func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response content type to plain text
	w.Header().Set("Content-Type", "application/plain")
	w.WriteHeader(http.StatusOK) // Send HTTP status 200
	fmt.Fprintln(w, "Test Github Actions") // Send response body
}

// AuthHandler returns the username from the request context as JSON
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get the username from the request context
	username := r.Context().Value("username")

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Send HTTP status 200

	// Create the response map
	resp := map[string]interface{}{
		"username": username,
	}

	// Encode the response as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

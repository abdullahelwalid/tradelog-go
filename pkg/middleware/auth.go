package middleware

import (
	"fmt"
	"strings"
	"net/http"
	"context"
	"encoding/json"

	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("AUTH ROUTE")

		// Check if the Authorization header is present and correctly formatted
		if authHeader == "" || len(strings.Split(authHeader, " ")) < 2 {
			fmt.Println("NO AUTH HEADER")
			// Set the response header to application/json
			w.Header().Set("Content-Type", "application/json")
			// Send a JSON error response
			w.WriteHeader(http.StatusUnauthorized) // 401
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Auth Header can't be empty",
			})
			return
		}

		// Initialize AWS config (same as in Fiber)
		auth, err := utils.InitAWSConfig()
		if err != nil {
			fmt.Println(err)
			// Set the response header to application/json
			w.Header().Set("Content-Type", "application/json")
			// Send a JSON error response
			w.WriteHeader(http.StatusInternalServerError) // 500
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Something went wrong",
			})
			return
		}

		// Validate the token (same as in Fiber)
		resp, err := auth.ValidateToken(strings.Split(authHeader, " ")[1])
		if err != nil {
			// Set the response header to application/json
			w.Header().Set("Content-Type", "application/json")
			// Send a JSON error response
			w.WriteHeader(http.StatusUnauthorized) // 401
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Something went wrong when validating token",
			})
			return
		}

		// Print the username (for debugging)
		fmt.Println(" **** USERNAME ****")
		fmt.Println(*resp.Username)

		// Set the username into the request context (for later use in handlers)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", *resp.Username)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

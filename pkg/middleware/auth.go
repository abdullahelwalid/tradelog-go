package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	authTypes "github.com/abdullahelwalid/tradelog-go/pkg/types"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
)

//middleware validates token, if token expired it get refreshed by checking if refresh token in cookies
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("AUTH ROUTE")

		if (authHeader == ""){

		authData, err := r.Cookie("authData")
		if (err != nil){
			switch {
			case errors.Is(err, http.ErrNoCookie):
				authHeader = ""	
				w.Header().Set("Content-Type", "application/json")
				// Send a JSON error response
				w.WriteHeader(http.StatusInternalServerError) // 500
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Something went wrong",
				})
			return

			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		}
		var authDataDeserialized authTypes.AuthCookies;
		err = json.Unmarshal([]byte(authData.Value), &authDataDeserialized)
		if err != nil {
			println(err)
			// Set the response header to application/json
			w.Header().Set("Content-Type", "application/json")
			// Send a JSON error response
			w.WriteHeader(http.StatusInternalServerError) // 500
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Something went wrong",
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
		resp, err := auth.ValidateToken(authDataDeserialized.AccessToken)
		if err != nil {
			// Try Refreshing Token
			resp, err := auth.RefreshToken(authDataDeserialized.RefreshToken, authDataDeserialized.Email)
			if err != nil {
				println(err)
				// Set the response header to application/json
				w.Header().Set("Content-Type", "application/json")
				// Send a JSON error response
				w.WriteHeader(http.StatusUnauthorized) // 401
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Refresh Token expired or invalid",
				})
				http.SetCookie(w, &http.Cookie{Name: "authData", Value: "", Expires: time.Unix(0, 0),})
				return
			}
			println(resp)
			//TODO - Add logic to update cookies
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
		return
	}


		// Check if the Authorization header is present and correctly formatted
		if authHeader == "" || len(strings.Split(authHeader, " ")) < 2 {	
			// Set the response header to application/json
			w.Header().Set("Content-Type", "application/json")
			// Send a JSON error response
			w.WriteHeader(http.StatusUnauthorized) // 401
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Refresh Token expired or invalid",
			})
			http.SetCookie(w, &http.Cookie{Name: "authData", Value: "", Expires: time.Unix(0, 0),})
			return
		} else {
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
		}
	})
}

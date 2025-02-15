package middleware


import (
	"net/http"
)


func MethodCheckMiddleware(next http.Handler, allowedMethods []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}


		// Check if the request method is allowed
		for _, method := range allowedMethods {
			if r.Method == method {
				// If method matches, continue to the next handler
				next.ServeHTTP(w, r)
				return
			}
		}
		
		// If method does not match any of the allowed methods, return 405 Method Not Allowed
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

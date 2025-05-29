package middleware

import (
	"log"
	"net/http"
)

// CORSMiddleware adds CORS headers to responses with proper configuration
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Define allowed origins
		allowedOrigins := []string{
			"http://localhost:5173", // Vite dev server
			"http://localhost:3000", // Alternative frontend port
			"http://localhost:8080", // Backend port (for debugging)
			"http://127.0.0.1:5173",
			"http://127.0.0.1:3000",
		}

		// Check if the origin is allowed
		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				originAllowed = true
				break
			}
		}

		// Set the appropriate CORS headers
		if originAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			// For requests without origin (like direct API calls) or unknown origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// Note: Cannot set credentials to true with wildcard origin
		}

		// Allow common HTTP methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Allow all headers the client might send
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With, Cache-Control, Origin")

		// Set max age for preflight requests
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			log.Printf("CORS preflight request from origin: %s", origin)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all requests with additional CORS debugging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			log.Printf("%s %s %s (Origin: %s)", r.RemoteAddr, r.Method, r.URL, origin)
		} else {
			log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		}
		next.ServeHTTP(w, r)
	})
}

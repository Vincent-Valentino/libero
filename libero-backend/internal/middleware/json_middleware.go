package middleware

import "net/http"

// JSONMiddleware ensures all responses have proper JSON headers and caching directives
func JSONMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set consistent JSON and caching headers
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Cache-Control", "private, must-revalidate")
        w.Header().Set("Vary", "Accept-Encoding")
        
        // Call the next handler
        next.ServeHTTP(w, r)
    })
}

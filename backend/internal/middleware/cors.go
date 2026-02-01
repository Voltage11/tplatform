package middleware

import (
	"fmt"
	"net/http"
)

// CORSMiddleware структура для CORS
type CORSMiddleware struct {
	allowedOrigins []string
}

// NewCORSMiddleware создает новый CORS middleware
func NewCORSMiddleware(allowedOrigins []string) *CORSMiddleware {
	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Handler основной обработчик CORS
func (c *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		requestID := r.Header.Get("X-Request-Id")

		allowed := false
		for _, allowedOrigin := range c.allowedOrigins {
			if allowedOrigin == "*" || origin == allowedOrigin {
				if allowedOrigin == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					// Не устанавливаем Allow-Credentials для wildcard
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				allowed = true
				break
			}
		}

		if !allowed && origin != "" {
			fmt.Printf("CORS [%s]: Origin %s NOT allowed\n", requestID, origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

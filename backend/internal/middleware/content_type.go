package middleware

import (
	"net/http"
	"strings"
)

func RequireJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем только методы, которые могут иметь тело
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			contentType := r.Header.Get("Content-Type")

			if !strings.HasPrefix(contentType, "application/json") {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("Content-Type must be application/json"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

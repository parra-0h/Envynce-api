package middleware

import (
	"net/http"

	"github.com/hans/config-service/internal/config"
)

func APIKeyAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-KEY")
			if key != cfg.APIKey {
				http.Error(w, "Unauthorized: Invalid API Key", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/hans/config-service/internal/domain"
	"github.com/hans/config-service/internal/service"
)

const (
	ContextKeyAPIKeyID   contextKey = "api_key_id"
	ContextKeyAPIKeyApps contextKey = "api_key_apps"
)

func APIKeyAuth(apiKeySvc *service.APIKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"Unauthorized: missing or invalid API Key"}`, http.StatusUnauthorized)
				return
			}

			plainKey := strings.TrimPrefix(authHeader, "Bearer ")
			apiKey, err := apiKeySvc.ValidateKey(r.Context(), plainKey)
			if err != nil {
				http.Error(w, `{"error":"Unauthorized: invalid or expired API Key"}`, http.StatusUnauthorized)
				return
			}

			// Store API Key info in context
			ctx := context.WithValue(r.Context(), ContextKeyAPIKeyID, apiKey.ID)
			ctx = context.WithValue(ctx, ContextKeyAPIKeyApps, apiKey.Applications)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAPIKeyID(ctx context.Context) uint {
	v, _ := ctx.Value(ContextKeyAPIKeyID).(uint)
	return v
}

func GetAPIKeyApps(ctx context.Context) []domain.Application {
	v, _ := ctx.Value(ContextKeyAPIKeyApps).([]domain.Application)
	return v
}

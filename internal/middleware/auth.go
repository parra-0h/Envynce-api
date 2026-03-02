package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
	ContextKeyName   contextKey = "name"
	ContextKeyEmail  contextKey = "email"
)

func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"Unauthorized: missing or invalid token"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"Unauthorized: invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"Unauthorized: invalid claims"}`, http.StatusUnauthorized)
				return
			}

			// Inject claims into context
			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims["user_id"])
			ctx = context.WithValue(ctx, ContextKeyRole, claims["role"])
			ctx = context.WithValue(ctx, ContextKeyName, claims["name"])
			ctx = context.WithValue(ctx, ContextKeyEmail, claims["email"])

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the context
func GetUserID(ctx context.Context) uint {
	v, _ := ctx.Value(ContextKeyUserID).(float64)
	return uint(v)
}

// GetUserRole extracts the role from the context
func GetUserRole(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyRole).(string)
	return v
}

// GetUserName extracts the name from the context
func GetUserName(ctx context.Context) string {
	v, _ := ctx.Value(ContextKeyName).(string)
	return v
}

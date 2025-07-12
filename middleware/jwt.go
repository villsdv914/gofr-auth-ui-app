package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

type contextKey string

const ClaimsKey contextKey = "jwt_claims"

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_secret_key"
	}
	return []byte(secret)
}

func JWTAuthMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Define public routes that don't need authentication
			publicRoutes := []string{
				"/signup",
				"/login",
				"/refresh-token",
				"/ui/",
			}

			// Check if current path is a public route
			for _, route := range publicRoutes {
				if strings.HasPrefix(r.URL.Path, route) {
					inner.ServeHTTP(w, r)
					return
				}
			}

			// Apply JWT validation for protected routes
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Missing authorization header", "path": "` + r.URL.Path + `"}`))
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Invalid authorization header format"}`))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return getJWTSecret(), nil
			})

			if err != nil || !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "Invalid token", "token_error": "` + err.Error() + `"}`))
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := context.WithValue(r.Context(), ClaimsKey, claims)
				r = r.WithContext(ctx)
			}

			inner.ServeHTTP(w, r)
		})
	}
}
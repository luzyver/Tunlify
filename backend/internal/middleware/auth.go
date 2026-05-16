package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	secret string
}

func NewAuth(secret string) *Auth {
	return &Auth{secret: secret}
}

func (a *Auth) Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(a.secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["sub"].(float64))
		username := claims["username"].(string)

		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

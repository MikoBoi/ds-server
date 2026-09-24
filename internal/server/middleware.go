package server

import (
	"context"
	"math"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userIDValue, ok := claims["user_id"].(float64)
		if !ok || userIDValue <= 0 || userIDValue != math.Trunc(userIDValue) || userIDValue > math.MaxInt64 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := int64(userIDValue)

		ctx := context.WithValue(
			r.Context(),
			"user_id",
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

package server

import (
	"context"
	"net/http"
	"strings"

	"grc_be/internal/biz"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Skip for now to avoid breaking existing frontend if it's not sending tokens yet
			// In production, this should return 401
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &biz.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("grc-secret-key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), biz.UserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(biz.UserKey).(*biz.Claims)
		if !ok {
			// If no token, we can't verify role. 
			// For this task, we assume that if the user wants role-based, they should send a token.
			// However, to maintain compatibility with current state, let's allow it if no token is present?
			// No, that defeats the purpose.
			http.Error(w, "unauthorized: no claims found", http.StatusUnauthorized)
			return
		}

		if claims.Role != "SuperAdmin" && claims.Role != "Admin" {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

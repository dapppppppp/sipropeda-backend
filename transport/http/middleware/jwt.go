package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"sipropeda-backend/transport/http/response"
)

// Samakan dengan secret di user.service.go
var jwtSecret = []byte("rahasia-sipropeda-skripsi")

// ContextKey untuk menyimpan ID User
type contextKey string
const UserIDKey contextKey = "user_id"

// JWTProtected adalah middleware untuk mengecek token
func JWTProtected(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akses ditolak: Token tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Akses ditolak: Token tidak valid atau kedaluwarsa"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Gagal membaca data token"})
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			response.WithJSON(w, http.StatusUnauthorized, map[string]string{"error": "User ID tidak valid"})
			return
		}

		userID, _ := uuid.FromString(userIDStr)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
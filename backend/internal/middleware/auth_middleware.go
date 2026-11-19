package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userRoleKey  contextKey = "user_role"
	userEmailKey contextKey = "user_email"
)

type TokenClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(
	jwtSecret string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Ambil token dari Authorization header.
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "Token tidak ditemukan", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Format token tidak valid", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				http.Error(w, "Token tidak valid", http.StatusUnauthorized)
				return
			}

			// Validasi JWT.
			claims := &TokenClaims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}

					return []byte(jwtSecret), nil
				},
				jwt.WithIssuer("driverent-api"),
				jwt.WithExpirationRequired(),
			)

			if err != nil || !token.Valid {
				http.Error(w, "Token tidak valid atau sudah kedaluwarsa", http.StatusUnauthorized)
				return
			}

			if claims.UserID <= 0 {
				http.Error(w, "Identitas user tidak valid", http.StatusUnauthorized)
				return
			}

			// Simpan identitas user ke request context.
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userRoleKey, claims.Role)
			ctx = context.WithValue(ctx, userEmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Mengambil User ID dari context.
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// Mengambil Role dari context.
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}

// Mengambil Email dari context.
func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}

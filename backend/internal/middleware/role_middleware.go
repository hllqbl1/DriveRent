package middleware

import "net/http"

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userRole, ok := GetUserRole(r.Context())
			if !ok {
				http.Error(w, "Role tidak ditemukan", http.StatusUnauthorized)
				return
			}

			// Periksa apakah role user diizinkan.
			allowed := false

			for _, role := range allowedRoles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Kamu tidak memiliki akses ke endpoint ini", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

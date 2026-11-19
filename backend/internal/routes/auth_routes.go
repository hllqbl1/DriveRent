package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/auth"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	handler auth.AuthHandler,
	jwtSecret string,
) {
	// Public routes
	mux.HandleFunc("POST /api/auth/register", handler.Register)
	mux.HandleFunc("POST /api/auth/login", handler.Login)

	// Protected route
	mux.Handle(
		"GET /api/auth/me",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(handler.GetProfile),
		),
	)
}

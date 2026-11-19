package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/middleware"
	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
)

func RegisterUserRoutes(
	mux *http.ServeMux,
	handler users.UserHandler,
	jwtSecret string,
) {
	auth := middleware.AuthMiddleware(jwtSecret)
	adminOnly := middleware.RequireRole("admin")

	mux.Handle("POST /api/users",
		auth(adminOnly(http.HandlerFunc(handler.CreateUser))),
	)

	mux.Handle("GET /api/users",
		auth(adminOnly(http.HandlerFunc(handler.GetAllUsers))),
	)

	mux.Handle("GET /api/users/{id}",
		auth(adminOnly(http.HandlerFunc(handler.GetUserByID))),
	)

	mux.Handle("PUT /api/users/{id}",
		auth(adminOnly(http.HandlerFunc(handler.UpdateUser))),
	)

	mux.Handle("DELETE /api/users/{id}",
		auth(adminOnly(http.HandlerFunc(handler.DeleteUser))),
	)
}

package routes

import (
	"net/http"

	"github.com/hllqbl1/DriveRent/backend/internal/modules/users"
)

func RegisterUserRoutes(
	mux *http.ServeMux,
	handler users.UserHandler,
) {
	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("GET /api/users", handler.GetAllUsers)
	mux.HandleFunc("GET /api/users/{id}", handler.GetUserByID)
	mux.HandleFunc("PUT /api/users/{id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", handler.DeleteUser)
}

package routes

import (
	"net/http"
	"time"

	"swiftly/backend/internal/app"
	"swiftly/backend/internal/pkg/response"
)

func Register(mux *http.ServeMux, a *app.App) {
	// 1. Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Welcome to Swiftly API", map[string]string{
			"time": time.Now().Format(time.RFC3339),
		})
	})

	// 2. Module User Routes
	a.User.Register(mux)

	// TODO: Register Store & Product routes here
}

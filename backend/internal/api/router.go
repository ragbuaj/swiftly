package api

import (
	"net/http"

	"swiftly/backend/internal/api/routes"
	"swiftly/backend/internal/app"
	"swiftly/backend/internal/middleware"
)

func NewRouter(a *app.App) http.Handler {
	mux := http.NewServeMux()

	// 1. Register all routes
	routes.Register(mux, a)

	// 2. Apply Global Middlewares
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.RateLimitMiddleware(10, 20)(handler)

	return handler
}

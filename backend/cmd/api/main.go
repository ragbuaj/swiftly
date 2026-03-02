package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/middleware"
	"swiftly/backend/internal/pkg/response"
	"swiftly/backend/internal/user/handler"
	"swiftly/backend/internal/user/repository"
	"swiftly/backend/internal/user/service"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Database (pgxpool)
	database.Init()
	pool := database.GetPool()
	defer database.Close()

	// Initialize User Module
	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, "Welcome to Swiftly API", map[string]string{
			"time": time.Now().Format(time.RFC3339),
		})
	})

	// User Routes
	userHandler.Register(mux)

	// Wrap mux with middleware from the internal package
	// Order: CORS -> Logging -> Mux
	handler := middleware.CORSMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	fmt.Printf("Swiftly Backend running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

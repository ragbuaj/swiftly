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
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/pkg/storage"
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

	// Initialize Redis
	database.InitRedis()
	defer database.CloseRedis()

	// Initialize Storage
	useSSL := os.Getenv("S3_USE_SSL") == "true"
	uploader, err := storage.NewMinioUploader(
		os.Getenv("S3_ENDPOINT"),
		os.Getenv("S3_ACCESS_KEY"),
		os.Getenv("S3_SECRET_KEY"),
		os.Getenv("S3_BUCKET_NAME"),
		os.Getenv("S3_PUBLIC_URL"),
		useSSL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO uploader: %v\n", err)
	}

	// Initialize Social Auth
	socialRegistry := socialauth.NewRegistry()

	// Initialize User Module
	userRepo := repository.NewUserRepository(pool)
	activityRepo := repository.NewActivityRepository(pool)
	userService := service.NewService(userRepo, activityRepo, uploader)
	userHandler := handler.NewUserHandler(userService, socialRegistry)

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
	// Order: Rate Limit -> CORS -> Logging -> Mux
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.RateLimitMiddleware(10, 20)(handler)

	fmt.Printf("Swiftly Backend running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}


package app

import (
	"log/slog"

	"swiftly/backend/internal/config"
	"swiftly/backend/internal/database"
	"swiftly/backend/internal/pkg/socialauth"
	"swiftly/backend/internal/pkg/storage"
	"swiftly/backend/internal/user/handler"
	"swiftly/backend/internal/user/repository"
	"swiftly/backend/internal/user/service"
)

type App struct {
	Config  *config.Config
	Storage storage.Uploader
	Social  *socialauth.Registry
	User    *handler.UserHandler
}

func Init(cfg *config.Config) *App {
	pool := database.GetPool()

	// Initialize Storage
	uploader, err := storage.NewMinioUploader(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.BucketName,
		cfg.S3.PublicURL,
		cfg.S3.UseSSL,
	)
	if err != nil {
		slog.Warn("Failed to initialize MinIO uploader", "error", err)
	}

	// Initialize Social Auth
	socialRegistry := socialauth.NewRegistry()

	// Initialize User Module
	userRepo := repository.NewUserRepository(pool)
	activityRepo := repository.NewActivityRepository(pool)
	userService := service.NewService(userRepo, activityRepo, uploader)
	userHandler := handler.NewUserHandler(userService, socialRegistry)

	return &App{
		Config:  cfg,
		Storage: uploader,
		Social:  socialRegistry,
		User:    userHandler,
	}
}

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"swiftly/backend/internal/api"
	"swiftly/backend/internal/app"
	"swiftly/backend/internal/config"
	"swiftly/backend/internal/database"
)

func main() {
	// 1. Setup Structured Logging (slog)
	setupLogger()

	// 2. Load Configuration
	cfg := config.Load()
	slog.Info("Configuration loaded successfully", "env", cfg.App.Environment)

	// 3. Setup Infrastructure
	database.Init()
	database.InitRedis()
	defer database.Close()
	defer database.CloseRedis()

	// 4. Initialize Modular App Container
	application := app.Init(cfg)

	// 5. Initialize Modular Router with Middlewares
	router := api.NewRouter(application)

	// 6. Start Server
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	slog.Info(fmt.Sprintf("Swiftly Backend running on http://localhost:%s", cfg.App.Port))
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func setupLogger() {
	var handler slog.Handler
	if os.Getenv("APP_ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil) // JSON logs for production
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}) // Text logs for development
	}
	slog.SetDefault(slog.New(handler))
}

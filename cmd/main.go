package main

import (
	"log/slog"
	"music-library/internal/config"
	"music-library/internal/db"
	"music-library/internal/handlers"
	"music-library/internal/repository"
	"music-library/internal/router"
	"music-library/internal/services"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Configuration loading error", "error", err)
		return
	}
	slog.Info("Configuration loaded successfully")

	database, err := db.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		slog.Error("Database connection failed", "host", cfg.DBHost, "port", cfg.DBPort, "error", err)
		return
	}
	defer database.Close()
	slog.Info("Database connection successfully")

	repo := repository.NewSongRepository(database, logger)
	service := services.NewSongService(repo, logger)
	handler := handlers.NewSongHandler(service, logger)

	r := router.NewRouter(handler)

	slog.Info("Starting server", "port", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		slog.Error("Server failed to start", "error", err)
		return
	}
}

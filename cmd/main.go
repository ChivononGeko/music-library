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

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Music Library API
// @version 1.0
// @description This is the API documentation for the Music Library
// @host localhost:8080
// @BasePath /api/v1
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

	err = runMigrations(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		slog.Error("Migrations failed", "error", err)
		return
	}
	slog.Info("Migrations executed successfully")

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

func runMigrations(host, port, user, password, dbname string) error {
	connString := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	m, err := migrate.New(
		"file://../migrations",
		connString,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

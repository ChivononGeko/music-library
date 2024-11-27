package main

import (
	"log"
	"music-library/internal/config"
	"music-library/internal/db"
	"music-library/internal/handlers"
	"music-library/internal/repository"
	"music-library/internal/router"
	"music-library/internal/services"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := db.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	repo := repository.NewSongRepository(db)
	service := services.NewSongService(repo, cfg.ExternalAPI)
	handler := handlers.NewSongHandler(service)

	r := router.NewRouter(handler)

	log.Printf("Starting server on port %s...", cfg.APIPort)
	log.Fatal(http.ListenAndServe(":"+cfg.APIPort, r))
}

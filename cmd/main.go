package main

import (
	"database/sql"
	"fmt"
	"log"
	"music-library/internal/config"
	"music-library/internal/handlers"
	"music-library/internal/repository"
	"music-library/internal/services"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	dataName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dataName)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	repo := repository.NewSongRepository(db)
	service := services.NewSongService(repo, cfg.ExternalAPI)
	handler := handlers.NewSongHandler(service)

	http.HandleFunc("/songs", handler.AddSongHandler)
	//Add more handler by mux

	log.Printf("Starting server on port %s...", cfg.APIPort)
	log.Fatal(http.ListenAndServe(":"+cfg.APIPort, nil))
}

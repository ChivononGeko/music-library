package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"music-library/internal/models"
	"net/http"
)

type SongService interface {
	AddSong(group, song string) error
	UpdateSong(id string, updateSong *models.Song) error
	GetAllSongs() ([]*models.Song, error)
	GetSong(id string) (*models.Song, error)
	DeleteSong(id string) error
}

type SongHandler struct {
	service SongService
	logger  *slog.Logger
}

func NewSongHandler(service SongService, logger *slog.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SongHandler) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.AddSong(request.Group, request.Song)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	song, err := h.service.GetSong(id)
	if err != nil {
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, song, http.StatusOK)
}

func (h *SongHandler) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := h.service.GetAllSongs()
	if err != nil {
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	sendSuccess(w, songs, http.StatusOK)
}

func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	var updateSong models.Song
	err := json.NewDecoder(r.Body).Decode(&updateSong)
	if err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")

	err = h.service.UpdateSong(id, &updateSong)
	if err != nil {
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := h.service.DeleteSong(id)
	if err != nil {
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

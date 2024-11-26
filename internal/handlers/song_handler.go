package handlers

import (
	"encoding/json"
	"fmt"
	"music-library/internal/models"
	"net/http"
)

type SongService interface {
	AddSong(group, song string) error
	UpdateSong(id string, group, song string) error
	GetAllSongs() ([]*models.Song, error)
	GetSong(id string) (*models.Song, error)
	DeleteSong(id string) error
}

type SongHandler struct {
	service SongService
}

func NewSongHandler(service SongService) *SongHandler {
	return &SongHandler{service: service}
}

func (h *SongHandler) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.AddSong(request.Group, request.Song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	song, err := h.service.GetSong(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(song); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *SongHandler) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	songs, err := h.service.GetAllSongs()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(songs); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")

	if err := h.service.UpdateSong(id, request.Group, request.Song); err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
}

func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if err := h.service.DeleteSong(id); err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"music-library/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// SongService interface for interacting with the song service.
type SongService interface {
	AddSong(group, song string) error
	UpdateSong(id string, updateSong *models.Song) error
	GetAllSongs() ([]*models.Song, error)
	GetSong(id string) (*models.Song, error)
	DeleteSong(id string) error
}

// SongHandler a handler for working with songs.
type SongHandler struct {
	service SongService
	logger  *slog.Logger
}

// NewSongHandler creates a new song handler
func NewSongHandler(service SongService, logger *slog.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}

// AddSongHandler adds a song.
// @Summary Add a song
// @Description Adds a new song to the library.
// @Tags songs
// @Accept json
// @Produce json
// @Param request body struct{ Group string `json:"group"`; Song string `json:"song"` } true "Song to add"
// @Success 201 {string} string "Successfully added"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /songs [post]
func (h *SongHandler) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received AddSong request")
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.logger.Error("Failed to decode AddSong request", "error", err)
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Adding song", "group", request.Group, "song", request.Song)
	err = h.service.AddSong(request.Group, request.Song)
	if err != nil {
		h.logger.Error("Failed to add song", "error", err.Error())
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Song added successfully", "group", request.Group, "song", request.Song)
	w.WriteHeader(http.StatusCreated)
}

// GetSongHandler gets information about the song.
// @Summary Get the song
// @Description Returns information about the song by its ID.
// @Tags songs
// @Produce json
// @Param id query string true "Song ID"
// @Success 200 {object} models.Song "Song information"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /songs [get]
func (h *SongHandler) GetSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received GetSong request", "id", id)

	song, err := h.service.GetSong(id)
	if err != nil {
		h.logger.Error("Failed to get song", "id", id, "error", err.Error())
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Song retrieved successfully", "id", id)
	sendSuccess(w, song, http.StatusOK)
}

// GetAllSongsHandler gets all the songs.
// @Summary Get all the songs
// @Description Returns a list of all the songs in the library.
// @Tags songs
// @Produce json
// @Success 200 {array} models.Song "List of songs"
// @Failure 500 {string} string "Server error"
// @Router /songs/all [get]
func (h *SongHandler) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received GetAllSongs request")

	songs, err := h.service.GetAllSongs()
	if err != nil {
		h.logger.Error("Failed to retrieve all songs", "error", err.Error())
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("All songs retrieved successfully", "count", len(songs))
	sendSuccess(w, songs, http.StatusOK)
}

// UpdateSongHandler updates the information about the song.
// @Summary Update Song
// @Description Updates information about an existing song by its ID.
// @Tags songs
// @Accept json
// @Param id query string true "Song ID"
// @Param song body models.Song true "Updated song information"
// @Success 204 "Successfully updated"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /songs [put]
func (h *SongHandler) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received UpdateSong request", "id", id)

	var updateSong models.Song
	err := json.NewDecoder(r.Body).Decode(&updateSong)
	if err != nil {
		h.logger.Error("Failed to decode UpdateSong request", "id", id, "error", err.Error())
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Updating song", "id", id, "song", updateSong)
	err = h.service.UpdateSong(id, &updateSong)
	if err != nil {
		h.logger.Error("Failed to update song", "id", id, "error", err.Error())
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Song updated successfully", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteSongHandler deletes the song.
// @Summary Delete the song
// @Description Deletes the song by its ID.
// @Tags songs
// @Param id query string true "Song ID"
// @Success 204 "Successfully deleted"
// @Failure 500 {string} string "Server error"
// @Router /songs [delete]
func (h *SongHandler) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received DeleteSong request", "id", id)

	err := h.service.DeleteSong(id)
	if err != nil {
		h.logger.Error("Failed to delete song", "id", id, "error", err.Error())
		sendError(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Song deleted successfully", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

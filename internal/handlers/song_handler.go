package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"music-library/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// SongService interface for interacting with the song service.
type SongService interface {
	AddSong(group, song string) error
	UpdateSong(id string, updateSong *models.Song) error
	GetAllSongs() ([]*models.Song, error)
	GetSong(id string) (*models.Song, error)
	DeleteSong(id string) error
	GetSongPaginated(filter map[string]string, page, pageSize int) ([]*models.Song, error)
	GetSongTextPaginated(id string, page, pageSize int) ([]string, error)
}

// SongHandler a handler for working with songs.
type SongHandler struct {
	service SongService
}

// NewSongHandler creates a new song handler
func NewSongHandler(service SongService) *SongHandler {
	return &SongHandler{
		service: service,
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
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.Error("Failed to decode AddSong request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Group == "" || request.Song == "" {
		slog.Error("Invalid Adding song request", "group", request.Group, "song", request.Song)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("Adding song", "group", request.Group, "song", request.Song)

	if err := h.service.AddSong(request.Group, request.Song); err != nil {
		slog.Error("Failed to add song", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Song added successfully", "group", request.Group, "song", request.Song)
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
	slog.Info("Received GetSong request", "id", id)

	song, err := h.service.GetSong(id)
	if err != nil {
		slog.Error("Failed to get song", "id", id, "error", err.Error())
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	slog.Info("Song retrieved successfully", "id", id)
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
	slog.Info("Received GetAllSongs request")

	songs, err := h.service.GetAllSongs()
	if err != nil {
		slog.Error("Failed to retrieve all songs", "error", err.Error())
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	slog.Info("All songs retrieved successfully", "count", len(songs))
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
	slog.Info("Received UpdateSong request", "id", id)

	var updateSong models.Song
	err := json.NewDecoder(r.Body).Decode(&updateSong)
	if err != nil {
		slog.Error("Failed to decode UpdateSong request", "id", id, "error", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slog.Info("Updating song", "id", id, "song", updateSong)
	err = h.service.UpdateSong(id, &updateSong)
	if err != nil {
		slog.Error("Failed to update song", "id", id, "error", err.Error())
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	slog.Info("Song updated successfully", "id", id)
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
	slog.Info("Received DeleteSong request", "id", id)

	err := h.service.DeleteSong(id)
	if err != nil {
		slog.Error("Failed to delete song", "id", id, "error", err.Error())
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}

	slog.Info("Song deleted successfully", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *SongHandler) GetSongPaginated(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := map[string]string{}

	// Читаем фильтры из query-параметров
	if group := query.Get("group"); group != "" {
		filter["group"] = group
	}
	if song := query.Get("song"); song != "" {
		filter["song"] = song
	}
	if text := query.Get("text"); text != "" {
		filter["text"] = text
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	slog.Info("Handling GetSongs request", "filter", filter, "page", page, "pageSize", pageSize)

	songs, err := h.service.GetSongPaginated(filter, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *SongHandler) GetSongTextPaginatedHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		http.Error(w, "Missing song ID", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 2
	}

	slog.Info("Handling GetSongTextPaginated request", "id", id, "page", page, "pageSize", pageSize)

	verses, err := h.service.GetSongTextPaginated(id, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

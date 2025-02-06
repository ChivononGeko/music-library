package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"music-library/internal/models"
	"net/http"
	"net/url"
)

type SongRepository interface {
	DeleteSongRepository(id string) error
	UpdateSongRepository(id string, song *models.Song) error
	GetAllSongsRepository() ([]*models.Song, error)
	GetSongRepository(id string) (*models.Song, error)
	AddSongRepository(song models.Song) error
	GetSongPaginated(filter map[string]string, page, pageSize int) ([]*models.Song, error)
	GetSongTextPaginated(id string, page, pageSize int) ([]string, error)
}

type SongService struct {
	repository SongRepository
	APIURL     string
}

func NewSongService(repository SongRepository) *SongService {
	return &SongService{
		repository: repository,
	}
}

func (s *SongService) AddSong(group, song string) error {
	groupEncoded := url.QueryEscape(group)
	songEncoded := url.QueryEscape(song)

	apiURL := fmt.Sprintf("%s?group=%s&song=%s", s.APIURL, groupEncoded, songEncoded)
	slog.Info("Fetching song details from API", "url", apiURL)

	resp, err := http.Get(apiURL)
	if err != nil {
		slog.Error("Failed to fetch song details from API", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("API returned non-OK status", "status", resp.StatusCode)
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read API response", "error", err)
		return err
	}

	var songDetail models.Song
	if err := json.Unmarshal(body, &songDetail); err != nil {
		slog.Error("Failed to unmarshal song data", "error", err)
		return err
	}

	fullSong, err := models.NewSong(group, song, songDetail.Text, songDetail.Link, songDetail.ReleaseDate)
	if err != nil {
		slog.Error("Error creating song model", "error", err)
		return err
	}

	if err := s.repository.AddSongRepository(*fullSong); err != nil {
		slog.Error("Failed to add song to repository", "song", fullSong, "error", err)
		return err
	}

	slog.Info("Successfully added song to repository", "song", fullSong)
	return nil
}

func (s *SongService) GetSong(id string) (*models.Song, error) {
	song, err := s.repository.GetSongRepository(id)
	if err != nil {
		slog.Error("Failed to get song from repository", "id", id, "error", err)
		return nil, err
	}

	slog.Info("Successfully fetched song from repository", "song", song)
	return song, nil
}

func (s *SongService) GetAllSongs() ([]*models.Song, error) {
	slog.Info("Fetching all songs from repository")

	songs, err := s.repository.GetAllSongsRepository()
	if err != nil {
		slog.Error("Failed to get all songs from repository", "error", err)
		return nil, err
	}

	slog.Info("Successfully fetched all songs", "count", len(songs))
	return songs, nil
}

func (s *SongService) UpdateSong(id string, updateSong *models.Song) error {
	slog.Info("Updating song in repository", "id", id, "song", updateSong)

	fullSong, err := models.NewSong(updateSong.GroupName, updateSong.SongName, updateSong.Text, updateSong.Link, updateSong.ReleaseDate)
	if err != nil {
		slog.Error("Error creating song model", "error", err)
		return err
	}

	if err := s.repository.UpdateSongRepository(id, fullSong); err != nil {
		slog.Error("Failed to update song in repository", "id", id, "error", err)
		return err
	}

	slog.Info("Successfully updated song", "id", id)
	return nil
}

func (s *SongService) DeleteSong(id string) error {
	slog.Info("Deleting song from repository", "id", id)

	if err := s.repository.DeleteSongRepository(id); err != nil {
		slog.Error("Failed to delete song from repository", "id", id, "error", err)
		return err
	}

	slog.Info("Successfully deleted song", "id", id)
	return nil
}

func (s *SongService) GetSongPaginated(filter map[string]string, page, pageSize int) ([]*models.Song, error) {
	slog.Info("Fetching filtered songs", "filter", filter, "page", page, "pageSize", pageSize)

	songs, err := s.repository.GetSongPaginated(filter, page, pageSize)
	if err != nil {
		slog.Error("Failed to fetch filtered songs", "error", err)
		return nil, fmt.Errorf("error fetching songs: %w", err)
	}

	slog.Info("Successfully fetched filtered songs", "count", len(songs))
	return songs, nil
}

func (s *SongService) GetSongTextPaginated(id string, page, pageSize int) ([]string, error) {
	slog.Info("Fetching song lyrics with pagination", "id", id, "page", page, "pageSize", pageSize)

	verses, err := s.repository.GetSongTextPaginated(id, page, pageSize)
	if err != nil {
		slog.Error("Failed to fetch song lyrics", "id", id, "error", err)
		return nil, fmt.Errorf("error fetching song lyrics: %w", err)
	}

	slog.Info("Successfully fetched song lyrics", "id", id, "verses_count", len(verses))
	return verses, nil
}

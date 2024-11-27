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
}

type SongService struct {
	repository SongRepository
	logger     *slog.Logger
	APIURL     string
}

func NewSongService(repository SongRepository, logger *slog.Logger) *SongService {
	return &SongService{
		repository: repository,
		logger:     logger,
	}
}

func (s *SongService) AddSong(group, song string) error {
	groupEncoded := url.QueryEscape(group)
	songEncoded := url.QueryEscape(song)

	url := fmt.Sprintf("%s?group=%s&song=%s", s.APIURL, groupEncoded, songEncoded)
	s.logger.Info("Fetching song details from API", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		s.logger.Error("Failed to fetch song details from API", "error", err)
		return fmt.Errorf("failed to fetch song details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("API returned non-OK status", "status", resp.StatusCode)
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read API response", "error", err)
		return fmt.Errorf("failed to read API response: %v", err)
	}

	var songDetail models.Song
	if err := json.Unmarshal(body, &songDetail); err != nil {
		s.logger.Error("Failed to unmarshal song data", "error", err)
		return fmt.Errorf("failed to unmarshal song data: %v", err)
	}

	fullSong, err := models.NewSong(group, song, songDetail.Text, songDetail.Link, songDetail.ReleaseDate)
	if err != nil {
		s.logger.Error("Error creating song model", "error", err)
		return fmt.Errorf("error creating song: %v", err)
	}

	if err := s.repository.AddSongRepository(*fullSong); err != nil {
		s.logger.Error("Failed to add song to repository", "song", fullSong, "error", err)
		return fmt.Errorf("failed to add song: %v", err)
	}

	s.logger.Info("Successfully added song to repository", "song", fullSong)
	return nil
}

func (s *SongService) GetSong(id string) (*models.Song, error) {
	s.logger.Info("Fetching song from repository", "id", id)

	song, err := s.repository.GetSongRepository(id)
	if err != nil {
		s.logger.Error("Failed to get song from repository", "id", id, "error", err)
		return nil, fmt.Errorf("error getting song from repository: %w", err)
	}

	s.logger.Info("Successfully fetched song from repository", "song", song)
	return song, nil
}

func (s *SongService) GetAllSongs() ([]*models.Song, error) {
	s.logger.Info("Fetching all songs from repository")

	songs, err := s.repository.GetAllSongsRepository()
	if err != nil {
		s.logger.Error("Failed to get all songs from repository", "error", err)
		return nil, fmt.Errorf("error getting songs from repository: %w", err)
	}

	s.logger.Info("Successfully fetched all songs", "count", len(songs))
	return songs, nil
}

func (s *SongService) UpdateSong(id string, updateSong *models.Song) error {
	s.logger.Info("Updating song in repository", "id", id, "song", updateSong)

	fullSong, err := models.NewSong(updateSong.GroupName, updateSong.SongName, updateSong.Text, updateSong.Link, updateSong.ReleaseDate)
	if err != nil {
		s.logger.Error("Error creating song model", "error", err)
		return fmt.Errorf("error creating song: %v", err)
	}

	if err := s.repository.UpdateSongRepository(id, fullSong); err != nil {
		s.logger.Error("Failed to update song in repository", "id", id, "error", err)
		return fmt.Errorf("failed to update song: %v", err)
	}

	s.logger.Info("Successfully updated song", "id", id)
	return nil
}

func (s *SongService) DeleteSong(id string) error {
	s.logger.Info("Deleting song from repository", "id", id)

	if err := s.repository.DeleteSongRepository(id); err != nil {
		s.logger.Error("Failed to delete song from repository", "id", id, "error", err)
		return fmt.Errorf("failed to delete song: %v", err)
	}

	s.logger.Info("Successfully deleted song", "id", id)
	return nil
}

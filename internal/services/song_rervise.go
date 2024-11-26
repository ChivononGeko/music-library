package services

import (
	"encoding/json"
	"fmt"
	"io"
	"music-library/internal/models"
	"net/http"
)

type SongRepository interface {
	DeleteSongRepository(id string) error
	UpdateSongRepository(id string, song models.Song) error
	GetAllSongsRepository() ([]*models.Song, error)
	GetSongRepository(id string) (*models.Song, error)
	AddSongRepository(song models.Song) error
}

type SongService struct {
	repository SongRepository
	APIURL     string
}

func NewSongService(repository SongRepository, APIURL string) *SongService {
	return &SongService{repository: repository, APIURL: APIURL}
}

func (s *SongService) AddSong(group, song string) error {
	url := fmt.Sprintf("%s?group=%s&song=%s", s.APIURL, group, song)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch song details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read API response: %v", err)
	}

	var songDetail models.Song
	if err := json.Unmarshal(body, &songDetail); err != nil {
		return fmt.Errorf("failed to unmarshal song data: %v", err)
	}

	fullSong := models.Song{
		GroupName:   group,
		SongName:    song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	return s.repository.AddSongRepository(fullSong)
}

func (s *SongService) GetSong(id string) (*models.Song, error) {
	song, err := s.repository.GetSongRepository(id)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}
	return song, nil
}

func (s *SongService) GetAllSongs() ([]*models.Song, error) {
	songs, err := s.repository.GetAllSongsRepository()
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}
	return songs, nil
}

func (s *SongService) UpdateSong(id string, song, group string) error {
	return nil
}

func (s *SongService) DeleteSong(id string) error {
	if err := s.repository.DeleteSongRepository(id); err != nil {
		return fmt.Errorf("service error: %w", err)
	}
	return nil
}

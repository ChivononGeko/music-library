package models

import (
	"fmt"
	"time"
)

type Song struct {
	ID          string `json:"id"`
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func NewSong(groupName, songName, text, link, releaseDate string) (*Song, error) {
	if groupName == "" || songName == "" {
		return nil, fmt.Errorf("group name and song name cannot be empty")
	}

	return &Song{
		ID:          generateID(),
		GroupName:   groupName,
		SongName:    songName,
		ReleaseDate: releaseDate,
		Text:        text,
		Link:        link,
	}, nil
}

func generateID() string {
	return fmt.Sprintf("song-%d", time.Now().UnixNano())
}

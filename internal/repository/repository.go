package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"music-library/internal/models"

	_ "github.com/lib/pq"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

func (r *SongRepository) Close() error {
	slog.Info("Closing database connection")
	return r.db.Close()
}

func (r *SongRepository) AddSongRepository(song models.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		slog.Error("Failed to add song", "group_name", song.GroupName, "song_name", song.SongName, "error", err)
		return fmt.Errorf("failed to add song: %w", err)
	}

	slog.Info("Song added successfully", "group_name", song.GroupName, "song_name", song.SongName)
	return nil
}

func (r *SongRepository) GetSongRepository(id string) (*models.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = $1`

	var song models.Song
	err := r.db.QueryRow(query, id).Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
	if err == sql.ErrNoRows {
		slog.Warn("No song found", "id", id)
		return nil, fmt.Errorf("no song found with id %s", id)
	} else if err != nil {
		slog.Error("Failed to execute query", "id", id, "error", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	slog.Info("Song retrieved successfully", "id", id)
	return &song, nil
}

func (r *SongRepository) GetAllSongsRepository() ([]*models.Song, error) {
	query := `SELECT id, group_name, song_name, release_date, text, link FROM songs`

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Failed to execute query for all songs", "error", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var songs []*models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			slog.Error("Failed to scan song row", "error", err)
			return nil, fmt.Errorf("failed to scan song row: %w", err)
		}
		songs = append(songs, &song)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over rows", "error", err)
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	slog.Info("Retrieved all songs successfully", "count", len(songs))
	return songs, nil
}

func (r *SongRepository) UpdateSongRepository(id string, song *models.Song) error {
	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`

	result, err := r.db.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		slog.Error("Failed to update song", "id", id, "error", err)
		return fmt.Errorf("failed to update song with id %s: %w", id, err)
	}

	if err := checkRowsAffected(result, id); err != nil {
		return err
	}

	slog.Info("Song updated successfully", "id", id)
	return nil
}

func (r *SongRepository) DeleteSongRepository(id string) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		slog.Error("Failed to delete song", "id", id, "error", err)
		return fmt.Errorf("failed to delete song with id %s: %w", id, err)
	}

	if err := checkRowsAffected(result, id); err != nil {
		return err
	}

	slog.Info("Song deleted successfully", "id", id)
	return nil
}

func checkRowsAffected(result sql.Result, id string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Failed to get rows affected", "id", id, "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		slog.Warn("No rows affected", "id", id)
		return fmt.Errorf("no song found with id %s", id)
	}

	slog.Info("Rows affected", "id", id, "rows_affected", rowsAffected)
	return nil
}

func (r *SongRepository) GetSongPaginated(filter map[string]string, page, pageSize int) ([]*models.Song, error) {
	query := `SELECT id, group_name, song_name, text, link, release_date 
	          FROM songs WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if group, ok := filter["group"]; ok {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", argID)
		args = append(args, "%"+group+"%")
		argID++
	}
	if song, ok := filter["song"]; ok {
		query += fmt.Sprintf(" AND song_name ILIKE $%d", argID)
		args = append(args, "%"+song+"%")
		argID++
	}
	if text, ok := filter["text"]; ok {
		query += fmt.Sprintf(" AND to_tsvector('russian', text) @@ plainto_tsquery('russian', $%d)", argID)
		args = append(args, text)
		argID++
	}

	query += fmt.Sprintf(" ORDER BY release_date DESC LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.Text, &song.Link, &song.ReleaseDate); err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}

	return songs, nil
}

func (r *SongRepository) GetSongTextPaginated(id string, page, pageSize int) ([]string, error) {
	query := `SELECT unnest(string_to_array(text, E'\n\n')) AS verse 
	          FROM songs WHERE id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, id, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []string
	for rows.Next() {
		var verse string
		if err := rows.Scan(&verse); err != nil {
			return nil, err
		}
		verses = append(verses, verse)
	}

	return verses, nil
}

package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Register sqlite3 driver

	"chlofisher.com/rosewood/internal/library"
)

type MusicStore struct {
	db *sql.DB
}

func NewMusicStore(conn *sql.DB) (*MusicStore) {
	return &MusicStore{db: conn}
}

func (s *MusicStore) Find(id string) (*library.Song, error) {
	var song library.Song
	err := s.db.QueryRow("SELECT id, path, title, artist, album FROM music WHERE id = ?", id).
		Scan(&song.ID, &song.Path, &song.Title, &song.Artist, &song.Album)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &song, err
}

func (s *MusicStore) Insert(song *library.Song) error {
	query := `
		INSERT INTO music (id, path, title, artist, album)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title = excluded.title,
			artist = excluded.artist,
			album = excluded.album;
	`

	_, err := s.db.Exec(query,
		song.ID,
		song.Path,
		song.Title,
		song.Artist,
		song.Album,
	)
	
	return err
}

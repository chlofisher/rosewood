package db

import (
	"fmt"
	"database/sql"

	"github.com/sqids/sqids-go"

	"chlofisher.com/rosewood/internal/metadata"
)

type MusicStore struct {
	db *sql.DB
	sqidGen *sqids.Sqids
}

func NewMusicStore(conn *sql.DB) (*MusicStore) {
	s, _ := sqids.New(sqids.Options{
		Alphabet: "VXL8ApI61Nw3MfnRgmuqDOU2KhFvTQ9zWoieClcH0bsPdGBjyZE7tkSarxJ4Y5",
		MinLength: 8,
	})
	return &MusicStore{db: conn, sqidGen: s}
}

func (s *MusicStore) Find(id string) (*metadata.Song, error) {
	var song metadata.Song
	err := s.db.QueryRow("SELECT id, path, title, artist, album FROM music WHERE id = ?", id).
		Scan(&song.ID, &song.Path, &song.Title, &song.Artist, &song.Album)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &song, err
}

func (s *MusicStore) Search(searchTerm string) ([]*metadata.Song, error) {
	var rows *sql.Rows
	var err error
	if searchTerm == "" {
		rows, err = s.db.Query("SELECT row_id, id, path, title, artist, album FROM music") 
	} else {
		query := `
			SELECT m.row_id, m.id, m.path, m.title, m.artist, m.album
			FROM music m
			JOIN music_search f ON m.row_id = f.rowid
			WHERE music_search MATCH ?
			ORDER BY bm25(music_search)
			LIMIT 50;
		`

		formattedSearch := searchTerm + "*"
		rows, err = s.db.Query(query, formattedSearch)
		if err != nil {
			return nil, fmt.Errorf("Search failed: %w", err)
		}
	}
	defer rows.Close()

	var results []*metadata.Song

	for rows.Next() {
		song := &metadata.Song{}
		err := rows.Scan(
			&song.Index,
			&song.ID,
			&song.Path,
			&song.Title,
			&song.Artist,
			&song.Album,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, song)
	}

	return results, nil
}

func (s *MusicStore) Insert(song *metadata.Song) error {
	query := `
		INSERT INTO music (path, title, artist, album)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			title = excluded.title,
			artist = excluded.artist,
			album = excluded.album
		RETURNING row_id;
	`

	var newIndex int64

	err := s.db.QueryRow(query,
		song.Path,
		song.Title,
		song.Artist,
		song.Album,
	).Scan(&newIndex)

	if err != nil {	
		return fmt.Errorf("Error inserting song into DB: %w", err)
	}

	song.Index = newIndex
	song.ID = s.generateHashID(newIndex) 

	_, err = s.db.Exec(`UPDATE music SET id = ? WHERE row_id = ?`, song.ID, song.Index)

	return err
}

func (s *MusicStore) generateHashID(index int64) string {
	id, _ := s.sqidGen.Encode([]uint64{uint64(index)})
	return id
}

package db

import (
	"fmt"
	"database/sql"
	_ "modernc.org/sqlite" // Register sqlite driver
)

func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		err = fmt.Errorf("Failed to open database connection: %w", err)
		return nil, err
	}

	if err := migrate(conn); err != nil {
		err = fmt.Errorf("Failed to migrate database: %w", err)
		return nil, err
	}

	return conn, nil
}

func migrate(conn *sql.DB) error {
	query := `
	-- 1. Create music table
	CREATE TABLE IF NOT EXISTS music (
		row_id INTEGER PRIMARY KEY AUTOINCREMENT, 
		id TEXT UNIQUE,
		path TEXT UNIQUE,
		title TEXT,
		artist TEXT,
		album TEXT
	);

	-- 2. Create FTS5 music search virtual table
	CREATE VIRTUAL TABLE IF NOT EXISTS music_search USING fts5(
		title,
		artist,
		album,
		content='music',
		content_rowid='row_id'
	);

	-- 3. Sync triggers
	CREATE TRIGGER IF NOT EXISTS music_ai AFTER INSERT ON music BEGIN
        INSERT INTO music_search(rowid, title, artist, album) 
        VALUES (new.row_id, new.title, new.artist, new.album);
   END;

   CREATE TRIGGER IF NOT EXISTS music_ad AFTER DELETE ON music BEGIN
       INSERT INTO music_search(music_search, rowid, title, artist, album) 
       VALUES('delete', old.row_id, old.title, old.artist, old.album);
   END;

   CREATE TRIGGER IF NOT EXISTS music_au AFTER UPDATE ON music BEGIN
       INSERT INTO music_search(music_search, rowid, title, artist, album) 
       VALUES('delete', old.row_id, old.title, old.artist, old.album);
       INSERT INTO music_search(rowid, title, artist, album) 
       VALUES (new.row_id, new.title, new.artist, new.album);
   END;
	`

	_, err := conn.Exec(query)
	return err
}

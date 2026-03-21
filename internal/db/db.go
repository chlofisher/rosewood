package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Register sqlite3 driver
)

func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func migrate(conn *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS music (
		id TEXT PRIMARY KEY,
		path TEXT UNIQUE,
		title TEXT,
		artist TEXT,
		album TEXT
	);`
	_, err := conn.Exec(query)
	return err
}

package main

import (
	"log"
	"net/http"
	"database/sql"

	"chlofisher.com/rosewood/internal/db"
	"chlofisher.com/rosewood/internal/api"
	"chlofisher.com/rosewood/internal/scanner"
)


func main() {
	port := ":8080"
	log.Printf("Server started at http://localhost%v/", port)

	// Open database connection
	conn, err := db.Open("media.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	musicStore := initMusicStore(conn, "/home/chloe/Music/")

	musicHandler := api.NewMusicHandler(musicStore)

	mux := http.NewServeMux()

	musicHandler.RegisterRoutes(mux)

	log.Fatal(http.ListenAndServe(port, mux))
}

func initMusicStore(conn *sql.DB, rootDir string) *db.MusicStore {
	musicStore := db.NewMusicStore(conn)

	scanner := scanner.FileScanner{Music: musicStore}

	err := scanner.Scan(rootDir) 
	if err != nil {
		log.Fatal("Error scanning music files: %v", err)
	}
	
	return musicStore
}

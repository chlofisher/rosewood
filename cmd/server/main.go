package main

import (
	"os"
	"path/filepath"
	"flag"
	"log"
	"strconv"
	"net/http"
	"database/sql"

	"chlofisher.com/rosewood/internal/db"
	"chlofisher.com/rosewood/internal/api"
	"chlofisher.com/rosewood/internal/scanner"
)

func init() {

}

func defaultDataDir() string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir != "" {
		return filepath.Join(dataDir, "rosewood")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not find user home directory!")
	}

	// Default to ~/.local/share if XDG_DATA_HOME doesn't exist
	return filepath.Join(home, ".local", "share", "rosewood")
}

func main() {
	var portNum int

	flag.IntVar(&portNum, "p", 8080, "The port the server will listen on")
	flag.IntVar(&portNum, "port", 8080, "The port the server will listen on")

	var dataPath string

	flag.StringVar(&dataPath, "d", defaultDataDir(), "The directory where the media database will be stored")
	flag.StringVar(&dataPath, "datadir", defaultDataDir(), "The directory where the media database will be stored")

	flag.Parse()

	port := ":" + strconv.Itoa(portNum)
	log.Printf("Server started at http://localhost%v/", port)

	// Open database connection
	err := os.MkdirAll(dataPath, 0755)
	dbPath := filepath.Join(dataPath, "media.db")

	log.Printf(dbPath)
	conn, err := db.Open(dbPath)
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

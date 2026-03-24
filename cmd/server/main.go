package main

import (
	"os"
	"path/filepath"
	"flag"
	"log"
	"strconv"
	"net/http"
	"database/sql"

	"chlofisher.com/rosewood/internal/server/db"
	"chlofisher.com/rosewood/internal/server/api"
	"chlofisher.com/rosewood/internal/server/scanner"
)

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

func parseFlags() (string, string, string) {
	var portNum int
	flag.IntVar(&portNum, "p", 8080, "The port the server will listen on")
	flag.IntVar(&portNum, "port", 8080, "The port the server will listen on")

	var dataPath string
	flag.StringVar(&dataPath, "data", defaultDataDir(), "The directory where the media database will be stored")

	var musicDir string
	flag.StringVar(&musicDir, "music", ".", "The directory where the media database will be stored")

	flag.Parse()

	if portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port: %d. Must be between 1 and 65536", portNum)
	}
	port := ":" + strconv.Itoa(portNum)

	musicDir, _ = filepath.Abs(musicDir)

	return port, dataPath, musicDir
}

func main() {
	port, dataPath, musicDir := parseFlags()

	log.Printf("Server started at http://localhost%v/", port)

	// Open database connection
	err := os.MkdirAll(dataPath, 0755)
	dbPath := filepath.Join(dataPath, "media.db")

	log.Printf("Opening DB at %s", dbPath)
	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	musicStore := initMusicStore(conn, musicDir)

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

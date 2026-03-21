package main

import (
	"fmt"
	"log"
	"net/http"

	"chlofisher.com/rosewood/internal/db"
	"chlofisher.com/rosewood/internal/streamer"
	"chlofisher.com/rosewood/internal/scanner"
)


func init() {
	fmt.Println("Initialising server...")
}


func main() {
	port := ":8080"
	fmt.Printf("Server started at http://localhost%v/\n", port)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Music Server is Online!")
	})

	conn, err := db.Open("media.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	musicStore := db.NewMusicStore(conn)
	musicHandler := streamer.NewMusicHandler(musicStore)
	scanner := scanner.FileScanner{Music: musicStore, RootDir: "/home/chloe/Music/"}
	if err := scanner.Scan(); err != nil {
		log.Fatal("Error scanning music files: %v", err)
	}

	http.Handle("/play/music", musicHandler)

	log.Fatal(http.ListenAndServe(port, nil))
}

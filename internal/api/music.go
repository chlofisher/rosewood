package api

import (
	"log"
	"encoding/json"
	"net/http"
	"chlofisher.com/rosewood/internal/db"
	"chlofisher.com/rosewood/internal/library"
)

type MusicHandler struct {
	Store *db.MusicStore
}

func NewMusicHandler(ms *db.MusicStore) *MusicHandler {
	return &MusicHandler {
		Store: ms,
	}	
}

func (h *MusicHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v0/songs/{id}/stream/", h.Stream)
	mux.HandleFunc("GET /api/v0/songs", h.Search)
}

func (h *MusicHandler) Stream(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	id := r.PathValue("id")
	
	log.Printf("Requested song ID %v", id)

	song, err := h.Store.Find(id)
	if err != nil {
		http.Error(w, "Song %v not found.", 404)
		return
	}

	log.Printf("Playing song: %s", song.Title)

	http.ServeFile(w, r, song.Path)
}

func (h *MusicHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	songs := []*library.Song{}
	songs, err := h.Store.Search(q)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		log.Printf("JSON encoding failed: %v", err)
	}
}

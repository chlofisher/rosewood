package api

import (
	"fmt"
	"net/http"
	"chlofisher.com/rosewood/internal/db"
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
	
	fmt.Printf("Requested song ID %v\n", id)

	song, err := h.Store.Find(id)
	if err != nil {
		http.Error(w, "Song %v not found.", 404)
		return
	}

	fmt.Printf("Playing song: %s\n", song.Title)

	http.ServeFile(w, r, song.Path)
}

func (h *MusicHandler) Search(w http.ResponseWriter, r *http.Request) {
	return
}

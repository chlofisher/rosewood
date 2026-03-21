package streamer

import (
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

func (h *MusicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	song, err := h.Store.Find(id)
	if err != nil {
		http.Error(w, "Song %v not found.", 404)
		return
	}

	stream(w, r, song.Path)
}

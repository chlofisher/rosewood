package streamer

import (
	"net/http"
)

func stream(w http.ResponseWriter, r *http.Request, path string) {
	http.ServeFile(w, r, path)
}

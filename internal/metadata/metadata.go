package metadata

type Song struct {
	Index int64 `json:"-"`
	ID string `json:"id"`
	Path string `json:"path"`
	Title string `json:"title"`
	Album string `json:"album"`
	Artist string `json:"artist"`
}

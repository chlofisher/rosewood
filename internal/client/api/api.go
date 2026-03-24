package api

import (
	"net/url"
)

func GetStreamEndpoint(server *url.URL, id string) *url.URL {
	return server.JoinPath("api", "v0", "songs", id, "stream")
}

func GetSearchEndpoint(server *url.URL, term string) *url.URL {
	u := server.JoinPath("api", "v0", "songs")
	params := url.Values{}
	params.Add("q", term)
	u.RawQuery = params.Encode()

	return u
}

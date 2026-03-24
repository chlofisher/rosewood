package client

import (
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"chlofisher.com/rosewood/internal/client/player"
	"chlofisher.com/rosewood/internal/client/api"
	"chlofisher.com/rosewood/internal/metadata"
)

type Client struct {
	Server *url.URL
	Player player.Player
}

func New(u *url.URL, p player.Player) (*Client, error) {
	return &Client{
		Server: u,
		Player: p,
	}, nil
}

func (c *Client) Close() {
	if c.Player != nil {
		c.Player.Close()
	}
}

func (c *Client) PlaySongID(id string) error {
	u := api.GetStreamEndpoint(c.Server, id)
	return c.Player.Play(u)	
}

func (c *Client) SearchSongs(term string) ([]*metadata.Song, error) {
	u := api.GetSearchEndpoint(c.Server, term)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %s", resp.Status)
	}

	var songs []*metadata.Song

	if err := json.NewDecoder(resp.Body).Decode(&songs); err != nil {
		return nil, err
	}

	return songs, nil
}

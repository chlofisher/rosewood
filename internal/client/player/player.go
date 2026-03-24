package player

import (
	"net/url"
	vlc "github.com/adrg/libvlc-go/v3"
)

type Player interface {
	Play(u *url.URL) error
	Close()
}

type LibVLC struct {
	instance *vlc.Player
	media *vlc.Media
}

func NewLibVLC() (*LibVLC, error) {
	err := vlc.Init(
		"--no-video",
		"--quiet",
		// "--file-caching=1000",
		// "--network-caching=1500",
		// "--clock-jitter=500",
	)
	if err != nil {
		return nil, err
	}

	player, err := vlc.NewPlayer()
	if err != nil {
		return nil, err
	}

	return &LibVLC{
		instance: player,
		media: nil,
	}, nil
}

func (v *LibVLC) Close() {
	if v.media != nil {
		v.media.Release()
	}

	if v.instance != nil {
		v.instance.Stop()
		v.instance.Release()
	}

	vlc.Release()	
}

func (v *LibVLC) Play(u *url.URL) error {
	v.instance.Stop()
	if v.media != nil {
		v.media.Release()
		v.media = nil
	}

	song, err := v.instance.LoadMediaFromURL(u.String())
	if err != nil {
		return err
	}

	v.media = song

	return v.instance.Play()
}

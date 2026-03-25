package player

import (
	"fmt"
	"errors"
	"net/url"
	vlc "github.com/adrg/libvlc-go/v3"
)

type LibVLC struct {
	instance *vlc.Player
	media *vlc.Media
}

func NewLibVLC() (*LibVLC, error) {
	err := vlc.Init(
		"--no-video",
		"--quiet",
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

func (v *LibVLC) Close() error {
    var errs []error

    if v.instance != nil {
        if v.instance.IsPlaying() {
            if err := v.instance.Stop(); err != nil {
                errs = append(errs, fmt.Errorf("vlc player stop: %w", err))
            }
        }
        if err := v.instance.Release(); err != nil {
            errs = append(errs, fmt.Errorf("vlc player release: %w", err))
        }
    }

    if v.media != nil {
        if err := v.media.Release(); err != nil {
            errs = append(errs, fmt.Errorf("vlc media release: %w", err))
        }
    }

    if err := vlc.Release(); err != nil {
        errs = append(errs, fmt.Errorf("vlc global release: %w", err))
    }

    return errors.Join(errs...)
}

func (v *LibVLC) Play(u *url.URL) error {
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

func (v *LibVLC) Pause() error {
	return v.instance.SetPause(true)
}

func (v *LibVLC) Resume() error {
	return v.instance.SetPause(false)
}

func (v *LibVLC) TogglePause() error {
	return v.instance.TogglePause()	
}

func (v *LibVLC) IsPaused() bool {
	state, err := v.instance.MediaState()
	if err != nil {
		return true
	}

	return state == vlc.MediaPaused
}

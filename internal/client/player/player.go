package player

import (
	"net/url"
)

type Player interface {
	Play(u *url.URL) error
	Pause() error
	Resume() error
	TogglePause() error
	Close() error
	IsPaused() bool
}

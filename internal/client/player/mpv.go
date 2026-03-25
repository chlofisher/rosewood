package player 

import (
	"net/url"
	mpv "github.com/gen2brain/go-mpv"
)

type LibMPV struct {
	instance *mpv.Mpv
}

func NewLibMPV() (*LibMPV, error) {
	m := mpv.New()
	err := m.Initialize()
	if err != nil {
		return nil, err
	}

	m.SetOptionString("vo", "null")
	m.SetOptionString("terminal", "no")

	return &LibMPV{
		instance: m,
	}, nil
}

func (m *LibMPV) Close() error {
	m.instance.TerminateDestroy()
	return nil
}

func (m *LibMPV) Play(u *url.URL) error {
	return m.instance.Command([]string{"loadfile", u.String(), "replace"})
}

func (m *LibMPV) Pause() error {
	return m.instance.SetProperty("pause", mpv.FormatFlag, true)
}

func (m *LibMPV) Resume() error {
	return m.instance.SetProperty("pause", mpv.FormatFlag, false)
}

func (m *LibMPV) TogglePause() error {
	return m.instance.Command([]string{"cycle", "pause"})
}

func (m *LibMPV) IsPaused() bool {
	p, _ := m.instance.GetProperty("pause", mpv.FormatFlag)
	return p.(bool)
}

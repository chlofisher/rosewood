package main

import (
	"fmt"
	"log"
	"errors"
	"flag"
	"net/url"
	tea "charm.land/bubbletea/v2"
	"chlofisher.com/rosewood/internal/client"
	"chlofisher.com/rosewood/internal/client/player"
	"chlofisher.com/rosewood/internal/client/tui"
)

func parseFlags() (string, string) {
	var server string
	flag.StringVar(&server, "server", "http://localhost:9000", "Server URL")
	flag.StringVar(&server, "s", "http://localhost:9000", "Server URL")

	var audio string
	flag.StringVar(&audio, "backend", "mpv", "Audio Backend")
	flag.StringVar(&audio, "b", "mpv", "Audio Backend")

	return server, audio
}

func initPlayer(s string) (player.Player, error) {
	var p player.Player
	var err error

	switch s {
	case "vlc":
		p, err = player.NewLibVLC()
	case "mpv":
		p, err = player.NewLibMPV()
	default:
		p = nil
		err = errors.New(fmt.Sprintf("invalid audio backend %s", s))
	}

	return p, err
}

func main() {
	server, audioBackend := parseFlags()

	u, _ := url.Parse(server)

	p, err := initPlayer(audioBackend)
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.New(u, p)	
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	app := tea.NewProgram(tui.New(c))
	if _, err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

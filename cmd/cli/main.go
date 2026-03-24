package main

import (
	"log"
	"net/url"
	tea "charm.land/bubbletea/v2"
	"chlofisher.com/rosewood/internal/client"
	"chlofisher.com/rosewood/internal/client/player"
	"chlofisher.com/rosewood/internal/client/tui"
)

func main() {
	u, _ := url.Parse("http://localhost:9000/")

	p, err := player.NewLibVLC()
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

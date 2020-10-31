package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// Gui represents a connected browser that shows status information
type Gui struct {
	Speed time.Duration
	c     *websocket.Conn
}

type guiMessage struct {
	Speed time.Duration `json:"speed"`
}

func (g *Gui) listenDuration() {
	var message guiMessage
	for g.c != nil {
		err := g.c.ReadJSON(&message)
		if err != nil {
			g.c = nil
		}
		g.Speed = message.Speed * time.Millisecond
	}
}

// NewGui creates a GUI instance from the specified websocket connection and duration
func NewGui(c *websocket.Conn, duration time.Duration) *Gui {
	g := &Gui{duration, c}
	go g.listenDuration()
	return g
}

// WriteStatus writes the current game status to the GUI, pausing for a specified amount of time
func (g *Gui) WriteStatus(status *GameStatus) {
	if g.c != nil {
		status.You = 0
		g.c.WriteJSON(status)
		time.Sleep(g.Speed)
	}
}

package main

import (
	"net/http"
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

func listenGuiSpeed(g *Gui) {
	var message guiMessage
	for g.c != nil {
		err := g.c.ReadJSON(&message)
		if err != nil {
			g.c = nil
		}
		g.Speed = message.Speed
	}
}

// HandleGuiSocket is a request handler that returns a GUI instance
func HandleGuiSocket(w http.ResponseWriter, r *http.Request) (*Gui, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	g := &Gui{200 * time.Millisecond, c}
	go listenGuiSpeed(g)
	return g, nil
}

// NewUnconnectedGui creates a GUI instance that represents the lack thereof
func NewUnconnectedGui() *Gui {
	return &Gui{0, nil}
}

// WriteStatus writes the current game status to the GUI, pausing for a specified amount of time
func (g *Gui) WriteStatus(status *Status) {
	if g.c != nil {
		status.You = 0
		gui.c.WriteJSON(status)
		time.Sleep(g.Speed)
	}
}

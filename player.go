package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

// Player represents a player connected to a websocket and their associated game information
type Player struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	Speed     int    `json:"speed"`
	Active    bool   `json:"active"`
	Name      string `json:"name"`
	conn      *websocket.Conn
	ch        chan string
}

// Input contains the action that a player took. It is sent to the server via the web socket.
type Input struct {
	Action string `json:"action"`
}

func (p *Player) createInputChannel() {
	ch := make(chan string)
	go func() {
		input := Input{}
		for {
			if p.conn != nil {
				err := p.conn.ReadJSON(&input)
				if err != nil {
					p.conn = nil
				}
				ch <- input.Action
			} else {
				ch <- "change_nothing"
			}
		}
	}()
	p.ch = ch
}

// ReadActionAndProcess reads an action from the player's web socket. If no action is sent within the deadline, the player drops out
func (p *Player) ReadActionAndProcess(id int, deadline time.Time, jump bool) {
	if p.conn != nil {

		var action string
		select {
		case action = <-p.ch:
		case <-time.After(deadline.Sub(time.Now().UTC())):
			p.Active = false
			return
		}
		if action != "turn_left" && action != "turn_right" && action != "slow_down" && action != "speed_up" {
			action = "change_nothing"
		}
		if action == "speed_up" {
			if p.Speed != 10 {
				p.Speed++
			}
		} else if action == "slow_down" {
			if p.Speed != 1 {
				p.Speed--
			}
		} else if action == "turn_left" {
			switch p.Direction {
			case "left":
				p.Direction = "down"
				break
			case "down":
				p.Direction = "right"
				break
			case "right":
				p.Direction = "up"
				break
			case "up":
				p.Direction = "left"
				break
			}
		} else if action == "turn_right" {
			switch p.Direction {
			case "left":
				p.Direction = "up"
				break
			case "down":
				p.Direction = "left"
				break
			case "right":
				p.Direction = "down"
				break
			case "up":
				p.Direction = "right"
				break
			}
		}

		for i := 1; i <= p.Speed; i++ {
			if p.Direction == "up" {
				p.Y--
			} else if p.Direction == "down" {
				p.Y++
			} else if p.Direction == "right" {
				p.X++
			} else if p.Direction == "left" {
				p.X--
			}

			if p.X >= config.Width || p.Y >= config.Height || p.X < 0 || p.Y < 0 {
				p.Active = false
				break
			}

			if !jump || i == 1 || i == p.Speed {
				if status.Cells[p.Y][p.X] == 0 {
					status.Cells[p.Y][p.X] = id
				} else {
					p.Active = false
					break
				}
			}
		}
	}
	log.Println("Player coordinates: ", id, p.X, p.Y)
}

// NewPlayer creates a new player that starts at the specified coordinates, with the specified websocket connection and name
func NewPlayer(x int, y int, c *websocket.Conn, name string) *Player {
	p := Player{Speed: 1, Active: true, Direction: DIRECTIONS[rand.Intn(4)], Name: name, X: x, Y: y, conn: c}
	p.createInputChannel()
	return &p
}

// CloseConnection closes the player's connection, if it exists
func (p *Player) CloseConnection() {
	if conn := p.conn; conn != nil {
		conn.Close()
	}
}

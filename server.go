package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const WIDTH = 40
const HEIGHT = 40
const NUMPLAYERS = 2

var DIRECTIONS = []string{"up", "right", "left", "down"}

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

type Status struct {
	Width    int                `json:"width"`
	Height   int                `json:"height"`
	Cells    [WIDTH][HEIGHT]int `json:"cells"`
	Players  map[int]*Player    `json:"players"`
	You      int                `json:"you"`
	Running  bool               `json:"running"`
	Deadline string             `json:"deadline"`
}

type Input struct {
	Action string `json:"action"`
}

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{CheckOrigin: checkOrigin}

var status Status

func speed(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	addPlayer(c)
}

func initGame() {
	status = Status{Width: WIDTH, Height: HEIGHT, Cells: [WIDTH][HEIGHT]int{}, Running: false, Players: make(map[int]*Player, 0)}
}

func addPlayer(c *websocket.Conn) {
	playerID := len(status.Players) + 1
	/* do not add player when a game is already running */
	if status.Running {
		return
	}
	status.Players[playerID] = &Player{Speed: 1, Active: true, Direction: DIRECTIONS[rand.Intn(4)], Name: strconv.Itoa(playerID), X: rand.Intn(40), Y: rand.Intn(40), conn: c}
	status.Cells[status.Players[playerID].Y][status.Players[playerID].X] = playerID
	/* start game if lobby is full now */
	if len(status.Players) == NUMPLAYERS {
		go game()
	}
}

func writeStatus() {
	/* send status */
	log.Print("writing status")
	for player := range status.Players {
		status.You = player
		// status.deadline
		if conn := status.Players[player].conn; conn != nil {
			err := conn.WriteJSON(&status)
			if err != nil {
				status.Players[player].conn = nil
			}
		}
	}
}

func processPlayers(deadline time.Time) {
	log.Print("reading actions")
	for playerID := range status.Players {
		processPlayer(playerID, deadline)
	}
}

func inputChannel(player *Player) chan string {
	ch := make(chan string)
	go func() {
		input := Input{}
		for {
			if player.conn != nil {
				err := player.conn.ReadJSON(&input)
				if err != nil {
					player.conn = nil
				}
				ch <- input.Action
			} else {
				ch <- "change_nothing"
			}
		}
	}()
	return ch
}

func processPlayer(playerID int, deadline time.Time) {
	if player := status.Players[playerID]; player.conn != nil {

		var action string
		select {
		case action = <-player.ch:
		case <-time.After(deadline.Sub(time.Now().UTC())):
			action = "change_nothing"
		}
		if action != "turn_left" && action != "turn_right" && action != "slow_down" && action != "speed_up" {
			action = "change_nothing"
		}
		if action == "speed_up" {
			if player.Speed != 10 {
				player.Speed++
			}
		} else if action == "slow_down" {
			if player.Speed != 1 {
				player.Speed--
			}
		} else if action == "turn_left" {
			switch player.Direction {
			case "left":
				player.Direction = "down"
				break
			case "down":
				player.Direction = "right"
				break
			case "right":
				player.Direction = "up"
				break
			case "up":
				player.Direction = "left"
				break
			}
		} else if action == "turn_right" {
			switch player.Direction {
			case "left":
				player.Direction = "up"
				break
			case "down":
				player.Direction = "left"
				break
			case "right":
				player.Direction = "down"
				break
			case "up":
				player.Direction = "right"
				break
			}
		}

		// TODO add jumping
		for i := 0; i < player.Speed; i++ {
			if player.Direction == "up" {
				player.Y--
			} else if player.Direction == "down" {
				player.Y++
			} else if player.Direction == "right" {
				player.X++
			} else if player.Direction == "left" {
				player.X--
			}

			if player.X >= WIDTH || player.Y >= HEIGHT || player.X < 0 || player.Y < 0 {
				player.Active = false
				break
			}

			if status.Cells[player.Y][player.X] == 0 {
				status.Cells[player.Y][player.X] = playerID
			} else {
				player.Active = false
				break
			}
		}
	}
}

func game() {
	if len(status.Players) == 0 {
		return
	}
	for playerID := range status.Players {
		status.Players[playerID].ch = inputChannel(status.Players[playerID])
	}
	status.Running = true
	for status.Running {
		timeout := time.Now().UTC().Add(time.Second * 10)
		status.Deadline = timeout.Format(time.RFC3339)
		writeStatus()
		processPlayers(timeout)

		/* receive actions */
		numLiving := 0
		for playerID := range status.Players {
			if status.Players[playerID].Active {
				numLiving++
			}
		}
		if numLiving < 2 {
			status.Running = false
			break
		}

		time.Sleep(time.Now().UTC().Sub(timeout))
	}
	writeStatus()

	/* close all connections */
	log.Print("closing connections")
	for player := range status.Players {
		if conn := status.Players[player].conn; conn != nil {
			conn.Close()
		}
	}
}

func main() {
	initGame()
	http.HandleFunc("/spe_ed", speed)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

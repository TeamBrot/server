package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// DIRECTIONS contains all possible directions
var DIRECTIONS = []string{"up", "right", "left", "down"}

// GameStatus c
type GameStatus struct {
	Width         int             `json:"width"`
	Height        int             `json:"height"`
	Cells         [][]int         `json:"cells"`
	Players       map[int]*Player `json:"players"`
	You           int             `json:"you"`
	Running       bool            `json:"running"`
	Deadline      string          `json:"deadline"`
	gui           *Gui
	config        *Config
	ocuppiedCells [][]map[int]int
}

func (s *GameStatus) checkPlayerConnections() bool {
	for _, player := range s.Players {
		if player.conn != nil {
			return true
		}
	}
	return false
}

func (s *GameStatus) processPlayers(deadline time.Time, jump bool) {
	log.Print("reading actions")
	var processedPlayers []int
	for playerID, player := range s.Players {
		processedPlayers = append(processedPlayers, playerID)
		player.ReadActionAndProcess(playerID, deadline, jump)
	}

	// check weither multiple players moved to the same field
	for _, playerID := range processedPlayers {
		for i := range s.ocuppiedCells[playerID] {
			for Y, X := range s.ocuppiedCells[playerID][i] {
				for _, otherPlayerID := range processedPlayers {
					for j := range s.ocuppiedCells[otherPlayerID] {
						for otherY, otherX := range s.ocuppiedCells[otherPlayerID][j] {
							if playerID != otherPlayerID {
								if Y == otherY && X == otherX {
									//deactivate players
									s.Players[playerID].Active = false
									s.Players[otherPlayerID].Active = false
									log.Print("Player ", playerID, " and player ", otherPlayerID, " moved both to field: ", Y, " ", X)
								}
							}
						}
					}
				}
				//log.Println("Occupied Cells: ", s.ocuppiedCells[playerID])
				cellValue := s.Cells[Y][X]
				if cellValue > 0 {
					s.Cells[Y][X] = cellValue - 10
				} else if cellValue == -11 {
					s.Cells[Y][X] = cellValue + 10
				}
			}
		}
	}
}

func (s *GameStatus) writeStatus() {
	for player := range s.Players {
		s.You = player
		// status.deadline
		if conn := s.Players[player].conn; conn != nil {
			err := conn.WriteJSON(s)
			if err != nil {
				s.Players[player].conn = nil
			}
		}
	}
	if s.gui != nil {
		s.gui.WriteStatus(s)
	}
}

func (s *GameStatus) closeConnections() {
	/* close all connections */
	log.Print("closing connections")
	for _, player := range s.Players {
		player.CloseConnection()
	}
}

func (s *GameStatus) getNumLiving() (int, string) {
	numLiving := 0
	var lastLiving string
	for _, player := range s.Players {
		if player.Active {
			numLiving++
			lastLiving = player.Name
		}
	}
	return numLiving, lastLiving
}

// AddPlayer adds a player to the current GameStatus. It closes the connection if the game is already running
func (s *GameStatus) AddPlayer(c *websocket.Conn, config *Config) {
	playerID := len(s.Players) + 1
	/* do not add player when a game is already running */
	if s.Running {
		log.Println("game already in progress, disconnecting new client")
		c.Close()
		return
	}
	s.Players[playerID] = NewPlayer(rand.Intn(s.config.Width), rand.Intn(s.config.Height), c, strconv.Itoa(playerID))
	s.Cells[s.Players[playerID].Y][s.Players[playerID].X] = playerID
}

// GetNumPlayers returns the amount of players inside the game
func (s *GameStatus) GetNumPlayers() int {
	return len(s.Players)
}

// GameLoop plays the game, reading and writing to the players' sockets. When it ends, it closes all connections and creates a new, empty game status
func (s *GameStatus) GameLoop() {
	if len(s.Players) <= 1 {
		return
	}
	s.Running = true
	turn := 1
	for s.Running {
		log.Println("Turn: ", turn)
		if !s.checkPlayerConnections() {
			s.Running = false
			log.Println("all connections closed, stopping game")
			break
		}
		timeout := time.Now().UTC().Add(time.Second * 1000)
		s.Deadline = timeout.Format(time.RFC3339)
		s.writeStatus()
		s.processPlayers(timeout, turn%6 == 0)

		numLiving, lastLiving := s.getNumLiving()
		/* receive actions */
		if numLiving < 2 {
			log.Println("all players but", lastLiving, "died, stopping game")
			s.Running = false
			break
		}

		time.Sleep(time.Now().UTC().Sub(timeout))
		turn++
	}
	s.writeStatus()
	s.closeConnections()

	/* the last game is over, we start a new one */
	sNew := NewGameStatus(s.config)
	sNew.gui = s.gui
	*s = *sNew
}

// NewGameStatus creates a new, non-running GameStatus with the specified configuration
func NewGameStatus(config *Config) *GameStatus {
	log.Println("initializing lobby")
	cells := make([][]int, config.Height)
	for i := range cells {
		cells[i] = make([]int, config.Width)
	}
	return &GameStatus{Width: config.Width, Height: config.Height, Cells: cells, Running: false, Players: make(map[int]*Player, 0), Deadline: "", You: 0, gui: nil, config: config, ocuppiedCells: make([][]map[int]int, config.Players+1)}
}

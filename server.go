package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ServerTime struct {
	Time         string `json:"time"`
	Milliseconds int    `json:"milliseconds"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

var status *GameStatus
var config Config

func speed(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	status.AddPlayer(c, &config)

	/* start game if lobby is full now */
	if status.GetNumPlayers() == config.Players {
		log.Println("starting game")
		go status.GameLoop()
	}
}

func speedGui(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "gui.html")
}

func speedTime(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now().UTC()
	serverTime := ServerTime{Time: timeNow.Format(time.RFC3339), Milliseconds: int(float64(timeNow.Nanosecond()) / 1000000.0)}
	js, err := json.Marshal(serverTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func speedGuiSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade gui", err)
		return
	}
	status.gui = NewGui(c, 200*time.Millisecond)
}

func main() {
	config = GetConfig()
	status = NewGameStatus(&config)

	http.HandleFunc("/spe_ed_time", speedTime)
	http.HandleFunc("/spe_ed", speed)
	http.HandleFunc("/spe_ed/gui", speedGuiSocket)
	http.HandleFunc("/", speedGui)
	log.Println("server started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

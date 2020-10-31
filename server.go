package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

	http.HandleFunc("/spe_ed", speed)
	http.HandleFunc("/spe_ed/gui", speedGuiSocket)
	http.HandleFunc("/", speedGui)
	log.Println("server started")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

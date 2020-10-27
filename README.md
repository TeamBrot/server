# spe_ed server
This project contains a program to emulate the official test environment found unter `wss://msoll.de/spe_ed`

## Libraries
This program uses the gorilla/websocket library that can be found under [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket).

To install the library, run `go get github.com/gorilla/websocket`. 

## Setup
To build the server, run `go build server.go`

To run the server, run `./server`

## TODO

1. Allow playing again when first match is over
2. Make global state local so that multiple games can be played at the same time
3. Randomize board width and height
4. Truly randomize everything (seed from timestamp)
5. Use environment variables for URL and MAXPLAYERS
6. Enhance logging
7. Add small trivial bots that play the game via WebSocket

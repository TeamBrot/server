# spe_ed server
This project contains a program to emulate the official test environment found unter `wss://msoll.de/spe_ed`

## Libraries
This program uses the gorilla/websocket library that can be found under [https://github.com/gorilla/websocket](https://github.com/gorilla/websocket).

To install the library, run `go get github.com/gorilla/websocket`. 

## Setup
To build the server, run `go build .`

To run the server, run `./server`

## TODO

[X] Allow playing again when first match is over
[ ] Make global state local so that multiple games can be played at the same time
[ ] Randomize board width and height
[ ] Truly randomize everything (seed from timestamp)
[ ] Use environment variables for URL and MAXPLAYERS
[X] Enhance logging
[X] Add small trivial bots that play the game via WebSocket

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections []*websocket.Conn
var gameState [9]string // Represents the Tic-Tac-Toe board

// WebSocket handler
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	connections = append(connections, conn)

	// Send the initial game state to the connected user
	err = conn.WriteJSON(gameState)
	if err != nil {
		log.Println(err)
		return
	}

	// Listen for incoming moves from the client
	for {
		var data map[string]interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println(err)
			break
		}

		// Handle the incoming game state update (board & turn)
		if board, ok := data["board"].([]interface{}); ok {
			// Update game state based on the board
			for i, v := range board {
				gameState[i] = v.(string)
			}
		}

		// Broadcast the updated game state to all clients
		for _, c := range connections {
			err := c.WriteJSON(map[string]interface{}{
				"board": gameState,
				"turn":  data["turn"],
			})
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", handleConnection)
	http.ListenAndServe(":8080", nil) // Listen on port 8080
}

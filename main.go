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

var connections = make(map[*websocket.Conn]string) // Map of connections to player roles
var gameState [9]string                            // Represents the Tic-Tac-Toe board
var currentPlayer = "X"                            // Tracks whose turn it is
var playerCount = 0                                // Tracks the number of connected players

// Broadcast the updated game state to all clients
func broadcastGameState() {
	for conn, player := range connections {
		err := conn.WriteJSON(map[string]interface{}{
			"board":  gameState,
			"turn":   currentPlayer,
			"player": player, // Send the assigned player role
		})
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(connections, conn)
		}
	}
}

// WebSocket handler
func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Assign a player role
	var assignedPlayer string
	if playerCount == 0 {
		assignedPlayer = "X"
	} else if playerCount == 1 {
		assignedPlayer = "O"
	} else {
		// Reject additional players
		conn.WriteJSON(map[string]interface{}{
			"error": "Game is full. Try again later.",
		})
		conn.Close()
		return
	}

	connections[conn] = assignedPlayer
	playerCount++
	log.Printf("Player %s connected", assignedPlayer)

	// Send the initial game state and player role
	err = conn.WriteJSON(map[string]interface{}{
		"board":  gameState,
		"turn":   currentPlayer,
		"player": assignedPlayer,
	})
	if err != nil {
		log.Println(err)
		return
	}

	// Handle incoming moves
	for {
		var data map[string]interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println(err)
			break
		}

		// Validate the move
		if turn, ok := data["player"].(string); ok && turn == currentPlayer {
			if board, ok := data["board"].([]interface{}); ok {
				for i, v := range board {
					gameState[i] = v.(string)
				}
				// Switch turn
				if currentPlayer == "X" {
					currentPlayer = "O"
				} else {
					currentPlayer = "X"
				}
				broadcastGameState()
			}
		} else {
			// Send an error message if it's not the player's turn
			conn.WriteJSON(map[string]interface{}{
				"error": "It's not your turn!",
			})
		}
	}

	// Remove the connection when the client disconnects
	delete(connections, conn)
	playerCount--
	log.Printf("Player %s disconnected", assignedPlayer)
}

func main() {
	http.HandleFunc("/", handleConnection)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

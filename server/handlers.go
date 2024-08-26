package server

import (
	"chess-game/game"
	"encoding/json"
	"fmt"
	"log"
)

func handleJoinGame(client *Client, data interface{}) {
	var joinData struct {
		PlayerName string `json:"playerName"`
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		sendError(client, "Error: Invalid data format for joining game.")
		return
	}

	err = json.Unmarshal(dataBytes, &joinData)
	if err != nil {
		sendError(client, "Error: Invalid data format for joining game.")
		return
	}

	if joinData.PlayerName == "" {
		sendError(client, "Error: Player name cannot be empty.")
		return
	}

	if client.Game != nil {
		sendError(client, "You are already in a game.")
		return
	}

	mutex.Lock()
	var gameToJoin *game.Game
	for c := range clients {
		if c.Game != nil && c.Game.Phase == "setup" && c != client && len(c.Game.Players) == 1 {
			gameToJoin = c.Game
			break
		}
	}

	if gameToJoin == nil {
		gameToJoin = game.NewGame("", "")
		gameToJoin.Players = []*game.Player{{}, {}}
	}

	client.Game = gameToJoin

	if client.Game.Players[0].Name == "" {
		client.PlayerID = 1
		client.Game.Players[0].Name = joinData.PlayerName
	} else {
		client.PlayerID = 2
		client.Game.Players[1].Name = joinData.PlayerName
	}
	mutex.Unlock()

	sendMessage(client, Message{Type: "game_setup", Data: client.PlayerID})
	log.Printf("Player %s joined as Player %d\n", joinData.PlayerName, client.PlayerID)

	if client.Game.Players[0].Name != "" && client.Game.Players[1].Name != "" {
		for c := range clients {
			if c.Game == client.Game {
				sendMessage(c, Message{Type: "game_start", Data: c.Game.GetGameState()})
			}
		}
	}
}

func handleSetupDone(client *Client, data interface{}) {
	placements, ok := data.([]string)
	if !ok || len(placements) != 5 {
		sendError(client, "Invalid character placements.")
		return
	}

	if client.PlayerID == 1 {
		client.Game.Players[0].Placement = placements
	} else if client.PlayerID == 2 {
		client.Game.Players[1].Placement = placements
	}

	if client.Game.Players[0].Placement != nil && client.Game.Players[1].Placement != nil {
		if err := client.Game.SetupGame(client.Game.Players[0].Placement, client.Game.Players[1].Placement); err != nil {
			sendError(client, fmt.Sprintf("Error setting up game: %v", err))
			return
		}

		for c := range clients {
			if c.Game == client.Game {
				sendMessage(c, Message{Type: "game_start", Data: c.Game.GetGameState()})
			}
		}
	}
}

func handleMakeMove(client *Client, data interface{}) {
	moveCmd, ok := data.(string)
	if !ok {
		sendError(client, "Invalid move command format")
		return
	}

	if client.Game.CurrentTurn != client.PlayerID-1 {
		sendError(client, "It's not your turn!")
		return
	}

	if client.Game.Phase != "playing" {
		sendError(client, "Cannot make a move during the setup phase.")
		return
	}

	currentPlayer := client.Game.GetCurrentPlayer()

	err := client.Game.MoveCharacter(currentPlayer, moveCmd)
	if err != nil {
		sendError(client, fmt.Sprintf("Invalid move: %s", err.Error()))
		return
	}

	for c := range clients {
		if c.Game == client.Game {
			sendGameState(c)
		}
	}

	if client.Game.IsGameOver() {
		winner := client.Game.Winner
		gameOverMessage := fmt.Sprintf("Game Over! %s wins!", winner.Name)

		for c := range clients {
			if c.Game == client.Game {
				sendMessage(c, Message{Type: "game_over", Data: gameOverMessage})
			}
		}
	}
}

func handleNewGame(client *Client) {
	if client.Game == nil {
		sendError(client, "You are not in a game.")
		return
	}

	client.Game.ResetGame()

	for c := range clients {
		if c.Game == client.Game {
			sendMessage(c, Message{Type: "game_setup", Data: c.PlayerID})
		}
	}
}

func sendGameState(client *Client) {
	gameState := client.Game.GetGameState()
	sendMessage(client, Message{Type: "game_state", Data: gameState})
}

func sendError(client *Client, message string) {
	sendMessage(client, Message{Type: "error", Data: message})
}

func sendMessage(client *Client, message Message) {
	err := client.Conn.WriteJSON(message)
	if err != nil {
		log.Println("Error sending message:", err)
		client.Conn.Close()
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
	}
}
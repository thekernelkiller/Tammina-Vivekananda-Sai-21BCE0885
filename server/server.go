package server

import (
	"chess-game/game"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ID       string
	Game     *game.Game
	PlayerID int
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	clients   = make(map[*Client]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	mutex = &sync.Mutex{}
)

func StartServer(port int) {
	http.HandleFunc("/", handleConnections)
	fmt.Printf("Server started on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		Conn: conn,
		ID:   generateUniqueID(),
	}

	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	go handleClient(client)
}

func generateUniqueID() string {
	return uuid.New().String()
}

func handleClient(client *Client) {
	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		client.Conn.Close()
		log.Printf("Client disconnected: %s", client.ID)
	}()

	for {
		var msg Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		switch msg.Type {
		case "join_game":
			handleJoinGame(client, msg.Data)
		case "setup_done": 
			handleSetupDone(client, msg.Data)
		case "make_move":
			handleMakeMove(client, msg.Data)
		case "new_game":
			handleNewGame(client)
		default:
			sendError(client, "Invalid message type")
		}
	}
}

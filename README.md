# Turn-based Chess-like Game (WebSocket-based Multiplayer)

## HitWicket Software Engineering Assignment Submission

### Description
Backend (Golang):
  - The server is implemented in Go, utilizing the "gorilla/websocket" library for websocket handling.
  - Game logic, including character movement, attack rules, and win conditions, is handled server-side to prevent cheating.
  - The server maintains the game state, validates moves, and broadcasts updates to connected clients.

Frontend (HTML, CSS, JavaScript):
  - A basic frontend provides a visual representation of the game board and characters.
  - JavaScript handles user input (move commands), communicates with the server via websockets, and updates the UI based on game state changes.

### Run instructions
1. Ensure Go compiler is installed.
2. Navigate to the project root directory in terminal and run `go mod tidy`.
3. Execute `go run main.go` to start the server.
4. Open `web/index.html` using Live Server in two different browsers to run game. 

package game

import (
	"fmt"
	"strings"
)

type Game struct {
	Board        *Board
	Players      []*Player
	CurrentTurn int
	Winner     *Player
	Phase        string
}

type GameState struct {
	Board        [5][5]string
	CurrentTurn int
	Winner      string
	Phase      string
}

func NewGame(player1Name, player2Name string) *Game {
	player1 := &Player{ID: "A", Name: player1Name}
	player2 := &Player{ID: "B", Name: player2Name}

	board := NewBoard()

	game := &Game{
		Board:        board,
		Players:      []*Player{player1, player2},
		CurrentTurn: 0,
		Winner:     nil,
		Phase:        "setup",
	}

	return game
}

func (g *Game) SetupGame(player1Placement, player2Placement []string) error {
	if len(player1Placement) != 5 || len(player2Placement) != 5 {
		return fmt.Errorf("invalid placement: each player must place 5 characters")
	}

	if err := g.placeCharacters(g.Players[0], player1Placement); err != nil {
		return err
	}

	if err := g.placeCharacters(g.Players[1], player2Placement); err != nil {
		return err
	}

	g.Phase = "playing"

	return nil
}

func (g *Game) placeCharacters(player *Player, placement []string) error {
	player.Characters = make([]*Character, 5)
	for i, charName := range placement {
		charType := charName[:2]

		var char *Character
		var row int
		if player.ID == "A" {
			row = 4
		} else {
			row = 0
		}

		switch charType {
		case "P1", "P2", "P3", "P4", "P5":
			char = NewPawn(player, charName, row, i)
		case "H1":
			char = NewHero1(player, charName, row, i)
		case "H2":
			char = NewHero2(player, charName, row, i)
		default:
			return fmt.Errorf("invalid character type: %s", charType)
		}

		player.Characters[i] = char
		g.Board.PlaceCharacter(char, row, i)
	}
	return nil
}

func (g *Game) MoveCharacter(player *Player, moveCmd string) error {
	if g.Phase != "playing" {
		return fmt.Errorf("cannot make a move during the setup phase")
	}

	parts := strings.Split(moveCmd, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid move command format")
	}

	charName := parts[0]
	moveDir := parts[1]

	var char *Character
	for _, c := range player.Characters {
		if c.Name == charName {
			char = c
			break
		}
	}
	if char == nil {
		return fmt.Errorf("character not found: %s", charName)
	}

	if err := char.Move(g.Board, moveDir); err != nil {
		return err
	}

	g.CurrentTurn = (g.CurrentTurn + 1) % 2

	return nil
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.Players[g.CurrentTurn]
}

func (g *Game) IsGameOver() bool {
	for _, player := range g.Players {
		charactersAlive := 0
		for _, char := range player.Characters {
			if char.Alive {
				charactersAlive++
			}
		}
		if charactersAlive == 0 {
			g.Winner = g.GetOpponent(player)
			return true
		}
	}

	return false
}

func (g *Game) GetOpponent(player *Player) *Player {
	for _, p := range g.Players {
		if p != player {
			return p
		}
	}
	return nil
}

func (g *Game) GetGameState() GameState {
	var boardRep [5][5]string
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			char := g.Board.Grid[row][col]
			if char == nil {
				boardRep[row][col] = ""
			} else {
				boardRep[row][col] = fmt.Sprintf("%s-%s", char.Player.ID, char.Name)
			}
		}
	}

	winner := ""
	if g.Winner != nil {
		winner = g.Winner.Name
	}

	return GameState{
		Board:        boardRep,
		CurrentTurn: g.CurrentTurn,
		Winner:      winner,
		Phase:      g.Phase,
	}
}

func (g *Game) ResetGame() {
	g.Board = NewBoard()
	g.CurrentTurn = 0
	g.Winner = nil
	g.Phase = "setup"

	for _, player := range g.Players {
		player.Placement = nil
	}
}
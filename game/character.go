package game

import (
	"fmt"
	"strings"
)

type CharacterType int

const (
	Pawn CharacterType = iota
	Hero1
	Hero2
)

type Character struct {
	Player *Player
	Name   string
	Type   CharacterType
	Row    int
	Col    int
	Alive  bool
}

func NewPawn(player *Player, name string, row, col int) *Character {
	return &Character{
		Player: player,
		Name:   name,
		Type:   Pawn,
		Row:    row,
		Col:    col,
		Alive:  true,
	}
}

func NewHero1(player *Player, name string, row, col int) *Character {
	return &Character{
		Player: player,
		Name:   name,
		Type:   Hero1,
		Row:    row,
		Col:    col,
		Alive:  true,
	}
}

func NewHero2(player *Player, name string, row, col int) *Character {
	return &Character{
		Player: player,
		Name:   name,
		Type:   Hero2,
		Row:    row,
		Col:    col,
		Alive:  true,
	}
}

func (c *Character) Move(board *Board, direction string) error {
	newRow, newCol, err := c.calculateMove(direction)
	if err != nil {
		return err
	}

	if !board.isValidPosition(newRow, newCol) {
		return fmt.Errorf("invalid move: out of bounds")
	}

	if c.Type == Hero1 || c.Type == Hero2 {
		if err := c.handleHeroMovement(board, newRow, newCol); err != nil {
			return err
		}
	} else {
		if !board.IsEmpty(newRow, newCol) {
			return fmt.Errorf("invalid move: target position occupied")
		}
	}

	return board.MoveCharacter(c, newRow, newCol)
}

func (c *Character) calculateMove(direction string) (newRow, newCol int, err error) {
	newRow = c.Row
	newCol = c.Col

	switch strings.ToUpper(direction) { 
	case "L":
		newCol--
	case "R":
		newCol++
	case "F":
		newRow-- 
	case "B":
		newRow++ 
	case "FL":
		newRow--
		newCol--
	case "FR":
		newRow--
		newCol++
	case "BL":
		newRow++
		newCol--
	case "BR":
		newRow++
		newCol++
	default:
		return newRow, newCol, fmt.Errorf("invalid move direction: %s", direction)
	}

	if c.Type == Hero1 || c.Type == Hero2 {
		newRow += (newRow - c.Row)
		newCol += (newCol - c.Col)
	}

	return newRow, newCol, nil
}

func (c *Character) handleHeroMovement(board *Board, newRow, newCol int) error {
	rowDiff := c.Row - newRow
	if rowDiff < 0 {
		rowDiff = -rowDiff
	}
	colDiff := c.Col - newCol
	if colDiff < 0 {
		colDiff = -colDiff
	}

	rowStep := 1
	if newRow < c.Row {
		rowStep = -1
	}
	colStep := 1
	if newCol < c.Col {
		colStep = -1
	}

	for i := 1; i <= rowDiff || i <= colDiff; i++ {
		row := c.Row + i*rowStep
		col := c.Col + i*colStep

		if !board.isValidPosition(row, col) {
			break
		}

		targetChar := board.GetCharacterAt(row, col)
		if targetChar != nil {
			targetChar.Alive = false
			board.Grid[row][col] = nil

			if targetChar.Player != c.Player {
				break
			}
		}
	}

	return nil
}
package game

import "fmt"

type Board struct {
	Grid [5][5]*Character
}

func NewBoard() *Board {
	board := &Board{}
	for row := 0; row < 5; row++ {
		for col := 0; col < 5; col++ {
			board.Grid[row][col] = nil
		}
	}
	return board
}

func (b *Board) PlaceCharacter(char *Character, row, col int) error {
	if !b.isValidPosition(row, col) {
		return fmt.Errorf("invalid board position: row %d, col %d", row, col)
	}

	if b.Grid[row][col] != nil {
		return fmt.Errorf("position already occupied")
	}

	b.Grid[row][col] = char
	char.Row = row
	char.Col = col
	return nil
}

func (b *Board) MoveCharacter(char *Character, newRow, newCol int) error {
	if !b.isValidPosition(newRow, newCol) {
		return fmt.Errorf("invalid move target position: row %d, col %d", newRow, newCol)
	}

	if b.Grid[newRow][newCol] != nil && b.Grid[newRow][newCol].Player == char.Player {
		return fmt.Errorf("cannot move to a position occupied by a friendly character")
	}

	b.Grid[char.Row][char.Col] = nil

	b.Grid[newRow][newCol] = char
	char.Row = newRow
	char.Col = newCol
	return nil
}

func (b *Board) IsEmpty(row, col int) bool {
	return b.Grid[row][col] == nil
}

func (b *Board) GetCharacterAt(row, col int) *Character {
	return b.Grid[row][col]
}

func (b *Board) isValidPosition(row, col int) bool {
	return row >= 0 && row < 5 && col >= 0 && col < 5
}
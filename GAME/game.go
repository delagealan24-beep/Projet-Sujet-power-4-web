package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	ROWS      = 6
	COLUMNS   = 7
	EMPTY     = 0
	PLAYER1   = 1
	PLAYER2   = 2
	WIN_COUNT = 4
)

type GameState int

const (
	Playing GameState = iota
	Player1Won
	Player2Won
	Draw
)

type Player struct {
	Name  string
	Color string
}

type Move struct {
	Player    int
	Column    int
	Row       int
	Timestamp time.Time
}

type Game struct {
	Grid          [ROWS][COLUMNS]int
	Players       [2]Player
	CurrentPlayer int
	State         GameState
	Moves         []Move
	StartTime     time.Time
	EndTime       time.Time
}

func NewGame(player1, player2 Player) *Game {
	game := &Game{
		Players:       [2]Player{player1, player2},
		CurrentPlayer: 0,
		State:         Playing,
		StartTime:     time.Now(),
		Moves:         make([]Move, 0),
	}

	for row := 0; row < ROWS; row++ {
		for col := 0; col < COLUMNS; col++ {
			game.Grid[row][col] = EMPTY
		}
	}

	return game
}

func (g *Game) Play(column int) error {
	if g.State != Playing {
		return errors.New("la partie est terminée")
	}

	if column < 0 || column >= COLUMNS {
		return errors.New("colonne invalide")
	}

	row := -1
	for r := ROWS - 1; r >= 0; r-- {
		if g.Grid[r][column] == EMPTY {
			row = r
			break
		}
	}

	if row == -1 {
		return errors.New("colonne pleine")
	}

	g.Grid[row][column] = g.CurrentPlayer + 1

	if g.checkWin(row, column) {
		if g.CurrentPlayer == 0 {
			g.State = Player1Won
		} else {
			g.State = Player2Won
		}
		g.EndTime = time.Now()
		return nil
	}

	// Vérifier l'égalité
	if g.isGridFull() {
		g.State = Draw
		g.EndTime = time.Now()
		return nil
	}
	g.CurrentPlayer = 1 - g.CurrentPlayer
	return nil
}
func (g *Game) checkWin(row, column int) bool {
	player := g.Grid[row][column]

	directions := [][2]int{
		{0, 1},  // horizontal
		{1, 0},  // vertical
		{1, 1},  // diagonale ↘
		{1, -1}, // diagonale ↙
	}

	for _, dir := range directions {
		count := 1
		count += g.countDirection(row, column, dir[0], dir[1], player)
		count += g.countDirection(row, column, -dir[0], -dir[1], player)

		if count >= WIN_COUNT {
			return true
		}
	}

	return false
}

func (g *Game) countDirection(startRow, startCol, deltaRow, deltaCol, player int) int {
	count := 0
	r, c := startRow+deltaRow, startCol+deltaCol

	for r >= 0 && r < ROWS && c >= 0 && c < COLUMNS && g.Grid[r][c] == player {
		count++
		r += deltaRow
		c += deltaCol
	}

	return count
}

func (g *Game) isGridFull() bool {
	for col := 0; col < COLUMNS; col++ {
		if g.Grid[0][col] == EMPTY {
			return false
		}
	}
	return true
}
func (g *Game) PrintGrid() {
	fmt.Println("\n  ┌─────────────┐")
	for row := 0; row < ROWS; row++ {
		fmt.Print("  │ ")
		for col := 0; col < COLUMNS; col++ {
			switch g.Grid[row][col] {
			case EMPTY:
				fmt.Print("○ ")
			case PLAYER1:
				fmt.Print("🔴")
			case PLAYER2:
				fmt.Print("🟡")
			}
		}
		fmt.Println("│")
	}
	fmt.Println("  └─────────────┘")
	fmt.Println("    1 2 3 4 5 6 7")
}

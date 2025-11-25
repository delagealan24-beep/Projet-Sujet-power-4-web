package main

import "time"

// Vérifie victoire à partir du dernier coup joué
func (g *Game) checkWin(row, column int) bool {
	player := g.Grid[row][column]
	directions := [][2]int{
		{0, 1},  // horizontal →
		{1, 0},  // vertical ↓
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
func (g *Game) PlayMove(column int) (int, int, error) {
	if g.State != Playing {
		return -1, -1, fmtError("la partie est terminée")
	}
	if column < 0 || column >= COLUMNS {
		return -1, -1, fmtError("colonne invalide")
	}
	// Trouver la ligne disponible
	row := -1
	for r := ROWS - 1; r >= 0; r-- {
		if g.Grid[r][column] == EMPTY {
			row = r
			break
		}
	}
	if row == -1 {
		return -1, -1, fmtError("colonne pleine")
	}
	player := g.CurrentPlayer + 1
	g.Grid[row][column] = player
	// Ajout du coup au journal
	g.Moves = append(g.Moves, Move{
		Player:    player,
		Column:    column,
		Row:       row,
		Timestamp: time.Now(),
	})
	// Vérification victoire
	if g.checkWin(row, column) {
		if g.CurrentPlayer == 0 {
			g.State = Player1Won
		} else {
			g.State = Player2Won
		}
		g.EndTime = time.Now()
		return row, player, nil
	}
	// Vérification égalité
	if g.isGridFull() {
		g.State = Draw
		g.EndTime = time.Now()
		return row, player, nil
	}
	// Changer de joueur
	g.CurrentPlayer = 1 - g.CurrentPlayer
	return row, player, nil
}

// Gestion d’erreur custom
type gameError struct{ s string }

func (e *gameError) Error() string { return e.s }
func fmtError(s string) error {
	return &gameError{s}
}

package main

import "fmt"

func main() {
	fmt.Println("🎮 POWER4 - Version Simplifiée")
	fmt.Println("===============================")

	player1 := Player{Name: "Alice", Color: "red"}
	player2 := Player{Name: "Bob", Color: "yellow"}

	game := NewGame(player1, player2)
	game.Play(3)
	game.Play(3)
	game.Play(2)

	game.PrintGrid()

	fmt.Printf("Joueur actuel: %s\n", game.Players[game.CurrentPlayer].Name)
	fmt.Printf("État du jeu: %v\n", game.State)
}

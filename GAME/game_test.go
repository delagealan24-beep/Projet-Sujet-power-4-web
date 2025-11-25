package main

import "testing"

func TestGame(t *testing.T) {
	player1 := Player{Name: "Test1", Color: "red"}
	player2 := Player{Name: "Test2", Color: "yellow"}

	game := NewGame(player1, player2)
	err := game.Play(0)
	if err != nil {
		t.Errorf("Erreur inattendue: %v", err)
	}
	for i := 0; i < 5; i++ {
		game.Play(0)
	}
	err = game.Play(0)
	if err == nil {
		t.Error("Devrait erreur sur colonne pleine")
	}
}

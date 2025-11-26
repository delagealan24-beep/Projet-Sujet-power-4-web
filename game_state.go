package main

import "time"

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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}
type Move struct {
	Player    int       `json:"player"`
	Column    int       `json:"column"`
	Row       int       `json:"row"`
	Timestamp time.Time `json:"timestamp"`
}

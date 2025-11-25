package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Game struct {
	ID            int                `json:"id"`
	Grid          [ROWS][COLUMNS]int `json:"grid"`
	Players       [2]Player          `json:"players"`
	CurrentPlayer int                `json:"current_player"`
	State         GameState          `json:"state"`
	Moves         []Move             `json:"moves"`
	StartTime     time.Time          `json:"start_time"`
	EndTime       time.Time          `json:"end_time,omitempty"`
	mu            sync.Mutex         `json:"-"`
}

func NewGame(p1, p2 Player) *Game {
	g := &Game{
		Players:       [2]Player{p1, p2},
		CurrentPlayer: 0,
		State:         Playing,
		StartTime:     time.Now(),
		Moves:         make([]Move, 0),
	}
	return g
}

var (
	games      = make(map[int]*Game)
	gamesMutex sync.RWMutex
	nextID     int32 = 0
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/games", gamesHandler)
	http.HandleFunc("/games/", singleGameHandler)
	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// POST /games  — créer une partie
// GET /games   — lister les parties
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload struct {
			Player1Name string `json:"player1_name"`
			Player2Name string `json:"player2_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "payload invalide", http.StatusBadRequest)
			return
		}
		p1 := Player{ID: 1, Name: payload.Player1Name}
		p2 := Player{ID: 2, Name: payload.Player2Name}
		g := NewGame(p1, p2)
		g.ID = int(atomic.AddInt32(&nextID, 1)) // commence à 1
		gamesMutex.Lock()
		games[g.ID] = g
		gamesMutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(g)
	case "GET":
		gamesMutex.RLock()
		list := make([]*Game, 0, len(games))
		for _, v := range games {
			list = append(list, v)
		}
		gamesMutex.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func singleGameHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/games/"):]
	parts := splitPath(path)
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "id invalide", http.StatusBadRequest)
		return
	}
	gamesMutex.RLock()
	g, ok := games[id]
	gamesMutex.RUnlock()
	if !ok {
		http.Error(w, "partie non trouvée", http.StatusNotFound)
		return
	}
	if len(parts) == 1 {
		if r.Method != "GET" {
			http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		json.NewEncoder(w).Encode(g)
		return
	}
	if len(parts) == 2 && parts[1] == "move" {
		if r.Method != "POST" {
			http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Column int `json:"column"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "payload invalide", http.StatusBadRequest)
			return
		}
		g.mu.Lock()
		_, _, err := g.PlayMove(payload.Column)
		g.mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(g)
		return
	}
	http.NotFound(w, r)
}

func splitPath(p string) []string {
	res := []string{}
	cur := ""
	for i := 0; i < len(p); i++ {
		if p[i] == '/' {
			if cur != "" {
				res = append(res, cur)
				cur = ""
			}
			continue
		}
		cur += string(p[i])
	}
	if cur != "" {
		res = append(res, cur)
	}
	return res
}

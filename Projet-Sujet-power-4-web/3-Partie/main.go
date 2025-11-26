package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GamePage struct {
	Grille        [6][7]string
	Tour          int
	JoueurActuel  string
	Joueur1       string
	Joueur2       string
	CouleurJ1     string
	CouleurJ2     string
	SymboleActuel string
	Message       string
	Winner        string
	Draw          bool
}

type GameState struct {
	Joueur1   string
	Joueur2   string
	CouleurJ1 string
	CouleurJ2 string
	Tour      int
	Grille    [6][7]string
	Started   bool
}

type ScoreEntry struct {
	Winner string
	Time   time.Time
}

var (
	templates  *template.Template
	game       GameState
	gameMutex  sync.Mutex
	scoreboard []ScoreEntry
	scoreMutex sync.Mutex
)

func main() {
	templates = template.Must(template.ParseGlob("templates/*.gohtml"))

	// static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/game/init", gameInitHandler)
	http.HandleFunc("/game/init/traitement", gameInitProcessHandler)
	http.HandleFunc("/game/play", gamePlayHandler)
	http.HandleFunc("/game/play/move", gameMoveHandler)
	http.HandleFunc("/game/end", gameEndHandler)
	http.HandleFunc("/game/scoreboard", scoreboardHandler)
	http.HandleFunc("/game/reset", resetHandler)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func render(w http.ResponseWriter, name string, data interface{}) {
	if err := templates.ExecuteTemplate(w, name+".gohtml", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "layout", map[string]interface{}{
		"Title":   "Power'4 Web",
		"Content": template.HTML(renderTemplateToString("homepage", nil)),
	})
}

func gameInitHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "layout", map[string]interface{}{
		"Title":   "Initialiser une partie",
		"Content": template.HTML(renderTemplateToString("game_init", nil)),
	})
}

func gameInitProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Impossible de lire le formulaire", http.StatusBadRequest)
		return
	}
	name1 := strings.TrimSpace(r.FormValue("name1"))
	name2 := strings.TrimSpace(r.FormValue("name2"))
	color1 := strings.TrimSpace(r.FormValue("color1"))
	color2 := strings.TrimSpace(r.FormValue("color2"))
	if name1 == "" {
		name1 = "Joueur 1"
	}
	if name2 == "" {
		name2 = "Joueur 2"
	}
	if color1 == "" {
		color1 = "🔴"
	}
	if color2 == "" {
		color2 = "🟡"
	}

	gameMutex.Lock()
	game = GameState{
		Joueur1:   name1,
		Joueur2:   name2,
		CouleurJ1: color1,
		CouleurJ2: color2,
		Tour:      1,
		Grille:    [6][7]string{},
		Started:   true,
	}
	gameMutex.Unlock()

	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

func gamePlayHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if !game.Started {
		http.Redirect(w, r, "/game/init", http.StatusSeeOther)
		return
	}

	var joueurActuel, symbole string
	if game.Tour%2 == 1 {
		joueurActuel = game.Joueur1
		symbole = game.CouleurJ1
	} else {
		joueurActuel = game.Joueur2
		symbole = game.CouleurJ2
	}

	page := GamePage{
		Grille:        game.Grille,
		Tour:          game.Tour,
		JoueurActuel:  joueurActuel,
		Joueur1:       game.Joueur1,
		Joueur2:       game.Joueur2,
		CouleurJ1:     game.CouleurJ1,
		CouleurJ2:     game.CouleurJ2,
		SymboleActuel: symbole,
		Message:       "C'est au tour de " + joueurActuel,
	}

	// vérifier victoire / égalité
	if win := checkVictory(game.Grille); win != "" {
		page.Winner = win
		page.Message = "Victoire de " + win + " !"
		// enregistrer score
		scoreMutex.Lock()
		scoreboard = append(scoreboard, ScoreEntry{Winner: win, Time: time.Now()})
		scoreMutex.Unlock()
	} else if isDraw(game.Grille) {
		page.Draw = true
		page.Message = "Match nul !"
	}

	render(w, "layout", map[string]interface{}{
		"Title":   "Partie en cours",
		"Content": template.HTML(renderTemplateToString("game_play", page)),
	})
}

func gameMoveHandler(w http.ResponseWriter, r *http.Request) {
	colStr := r.URL.Query().Get("col")
	col, err := strconv.Atoi(colStr)
	if err != nil || col < 0 || col > 6 {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}

	gameMutex.Lock()
	defer gameMutex.Unlock()

	if !game.Started {
		http.Redirect(w, r, "/game/init", http.StatusSeeOther)
		return
	}

	// si partie déjà terminée, rediriger vers /game/end
	if checkVictory(game.Grille) != "" || isDraw(game.Grille) {
		http.Redirect(w, r, "/game/end", http.StatusSeeOther)
		return
	}

	var symbole string
	if game.Tour%2 == 1 {
		symbole = game.CouleurJ1
	} else {
		symbole = game.CouleurJ2
	}

	if placeToken(&game.Grille, col, symbole) {
		game.Tour++
	}

	http.Redirect(w, r, "/game/play", http.StatusSeeOther)
}

func gameEndHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	winner := checkVictory(game.Grille)
	draw := isDraw(game.Grille)

	page := GamePage{
		Grille:  game.Grille,
		Tour:    game.Tour,
		Joueur1: game.Joueur1,
		Joueur2: game.Joueur2,
	}

	if winner != "" {
		page.Winner = winner
		page.Message = "Victoire de " + winner + " !"
	} else if draw {
		page.Draw = true
		page.Message = "Match nul !"
	} else {
		page.Message = "Partie non terminée."
	}

	render(w, "layout", map[string]interface{}{
		"Title":   "Fin de partie",
		"Content": template.HTML(renderTemplateToString("game_end", page)),
	})
}

func scoreboardHandler(w http.ResponseWriter, r *http.Request) {
	scoreMutex.Lock()
	defer scoreMutex.Unlock()
	render(w, "layout", map[string]interface{}{
		"Title":   "Classement",
		"Content": template.HTML(renderTemplateToString("scoreboard", scoreboard)),
	})
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	game = GameState{}
	gameMutex.Unlock()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderTemplateToString(name string, data interface{}) string {
	var sb strings.Builder
	if err := templates.ExecuteTemplate(&sb, name+".gohtml", data); err != nil {
		log.Println("Erreur render:", err)
		return ""
	}
	return sb.String()
}

func placeToken(grid *[6][7]string, col int, symbol string) bool {
	if col < 0 || col > 6 {
		return false
	}
	for row := 5; row >= 0; row-- {
		if grid[row][col] == "" {
			grid[row][col] = symbol
			return true
		}
	}
	return false
}

func isDraw(grid [6][7]string) bool {
	for c := 0; c < 7; c++ {
		if grid[0][c] == "" {
			return false
		}
	}
	// si pas de gagnant et première ligne pleine => match nul
	if checkVictory(grid) == "" {
		return true
	}
	return false
}

func checkVictory(grid [6][7]string) string {
	// check every cell; if not empty, check 4 directions
	dirs := [][2]int{
		{0, 1},  // horiz
		{1, 0},  // vert
		{1, 1},  // diag down-right
		{1, -1}, // diag down-left
	}
	for r := 0; r < 6; r++ {
		for c := 0; c < 7; c++ {
			s := grid[r][c]
			if s == "" {
				continue
			}
			for _, d := range dirs {
				count := 1
				nr, nc := r, c
				for i := 0; i < 3; i++ {
					nr += d[0]
					nc += d[1]
					if nr < 0 || nr >= 6 || nc < 0 || nc >= 7 {
						break
					}
					if grid[nr][nc] == s {
						count++
					} else {
						break
					}
				}
				if count >= 4 {
					// return player's name based on symbol
					if s == game.CouleurJ1 {
						return game.Joueur1
					} else if s == game.CouleurJ2 {
						return game.Joueur2
					}
					return s
				}
			}
		}
	}
	return ""
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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
}

type GameState struct {
	Joueur1   string
	Joueur2   string
	CouleurJ1 string
	CouleurJ2 string
	Tour      int
	Grille    [6][7]string
}

var (
	currentGame GameState
	templates   *template.Template
	mu          sync.Mutex
)

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if err := templates.ExecuteTemplate(w, name+".html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func placerJeton(col int, symbole string) bool {
	if col < 0 || col > 6 {
		return false
	}
	for row := 5; row >= 0; row-- {
		if currentGame.Grille[row][col] == "" {
			currentGame.Grille[row][col] = symbole
			return true
		}
	}
	return false
}

func main() {
	templates = template.Must(template.ParseGlob("./templates/*.html"))

	rootDoc, _ := os.Getwd()
	fs := http.FileServer(http.Dir(rootDoc + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "Homepage", nil)
	})

	http.HandleFunc("/game/init", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "GameInit", nil)
	})

	http.HandleFunc("/game/init/traitement", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		mu.Lock()
		currentGame = GameState{
			Joueur1:   r.FormValue("name"),
			Joueur2:   r.FormValue("name2"),
			CouleurJ1: r.FormValue("color1"),
			CouleurJ2: r.FormValue("color2"),
			Tour:      1,
			Grille:    [6][7]string{},
		}
		mu.Unlock()

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	})

	http.HandleFunc("/game/play", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		var joueurActuel, symbole string
		if currentGame.Tour%2 == 1 {
			joueurActuel = currentGame.Joueur1
			symbole = currentGame.CouleurJ1
		} else {
			joueurActuel = currentGame.Joueur2
			symbole = currentGame.CouleurJ2
		}
		data := GamePage{
			Grille:        currentGame.Grille,
			Tour:          currentGame.Tour,
			JoueurActuel:  joueurActuel,
			Joueur1:       currentGame.Joueur1,
			Joueur2:       currentGame.Joueur2,
			CouleurJ1:     currentGame.CouleurJ1,
			CouleurJ2:     currentGame.CouleurJ2,
			SymboleActuel: symbole,
			Message:       "C'est au tour de " + joueurActuel + " de jouer !",
		}
		mu.Unlock()

		renderTemplate(w, "GamePlay", data)
	})

	http.HandleFunc("/game/play/move", func(w http.ResponseWriter, r *http.Request) {
		colStr := r.URL.Query().Get("col")
		col, err := strconv.Atoi(colStr)
		if err != nil {
			http.Error(w, "Colonne invalide", http.StatusBadRequest)
			return
		}

		mu.Lock()
		var symbole string
		if currentGame.Tour%2 == 1 {
			symbole = currentGame.CouleurJ1
		} else {
			symbole = currentGame.CouleurJ2
		}

		if placerJeton(col, symbole) {
			currentGame.Tour++
		}
		mu.Unlock()

		http.Redirect(w, r, "/game/play", http.StatusSeeOther)
	})

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

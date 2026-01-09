package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Joueur struct {
	Pseudo  string
	Couleur string
}

type Game struct {
	Grille       [6][7]string
	Joueur1      Joueur
	Joueur2      Joueur
	Message      string
	JoueurActuel string
	Tour         int
	Fin          bool
	Gagnant      Joueur
	Egalite      bool
	VictoryType  string
}

type ScoreboardData struct {
	Scores     map[string]int
	Historique []string
}

func lireHistorique() []string {
	var historique []string
	file, err := os.Open("../text.txt")
	if err != nil {
		return historique
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		historique = append(historique, scanner.Text())
	}
	return historique
}

func enregistrerHistorique(j1, j2, vainqueur string) {
	file, err := os.OpenFile("../text.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Erreur historique :", err)
		return
	}
	defer file.Close()

	date := time.Now().Format("02/01/2006")
	var ligne string
	if vainqueur == "" {
		ligne = fmt.Sprintf("%s | %s vs %s | Égalité\n", date, j1, j2)
	} else {
		ligne = fmt.Sprintf("%s | %s vs %s | Vainqueur : %s\n", date, j1, j2, vainqueur)
	}
	file.WriteString(ligne)
}

var game Game

func initialiserGrille() [6][7]string {
	var g [6][7]string
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			g[i][j] = " "
		}
	}
	return g
}
func afficherGrille(grille [6][7]string) {
	fmt.Println("État de la grille :")
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if grille[i][j] == " " {
				fmt.Print("| |")
			} else {
				fmt.Printf("|%s|", grille[i][j][:1])
			}
		}
		fmt.Println()
	}
	fmt.Println(" 0  1  2  3  4  5  6")
	fmt.Println()
}

func switchJoueur(current string) string {
	if current == game.Joueur1.Couleur {
		return game.Joueur2.Couleur
	}
	return game.Joueur1.Couleur
}

func ajouterPion(grille *[6][7]string, colonne int, couleur string) bool {
	if colonne < 0 || colonne >= 7 {
		return false
	}
	for i := 5; i >= 0; i-- {
		if grille[i][colonne] == " " {
			grille[i][colonne] = couleur
			return true
		}
	}
	return false
}

func verifierVictoireType(grille [6][7]string, couleur string) string {

	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if grille[i][j] == couleur && grille[i][j+1] == couleur &&
				grille[i][j+2] == couleur && grille[i][j+3] == couleur {
				return "horizontal"
			}
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 7; j++ {
			if grille[i][j] == couleur && grille[i+1][j] == couleur &&
				grille[i+2][j] == couleur && grille[i+3][j] == couleur {
				return "vertical"
			}
		}
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if grille[i][j] == couleur && grille[i+1][j+1] == couleur &&
				grille[i+2][j+2] == couleur && grille[i+3][j+3] == couleur {
				return "diag_desc"
			}
		}
	}

	for i := 3; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if grille[i][j] == couleur && grille[i-1][j+1] == couleur &&
				grille[i-2][j+2] == couleur && grille[i-3][j+3] == couleur {
				return "diag_asc"
			}
		}
	}
	return ""
}

func verifierEgalite(grille [6][7]string) bool {
	for i := 0; i < 6; i++ {
		for j := 0; j < 7; j++ {
			if grille[i][j] == " " {
				return false
			}
		}
	}
	return true
}

func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func emptyOrValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "\u00A0"
	}
	return s
}

func lireScores() map[string]int {
	scores := make(map[string]int)
	file, err := os.Open("../text.txt")
	if err != nil {
		return scores
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ligne := scanner.Text()
		parts := strings.Split(ligne, ":")
		if len(parts) == 2 {
			score, err := strconv.Atoi(parts[1])
			if err == nil {
				scores[parts[0]] = score
			}
		}
	}
	return scores
}

func enregistrerVictoire(pseudo string) {
	scores := lireScores()
	scores[pseudo]++
	file, err := os.Create("../text.txt")
	if err != nil {
		fmt.Println("Erreur d’écriture :", err)
		return
	}
	defer file.Close()

	for joueur, score := range scores {
		fmt.Fprintf(file, "%s:%d\n", joueur, score)
	}
}

func main() {

	fmt.Println("DEBUG: registering template functions: seq, contains, emptyOrValue")

	listeTemplate, errTemplate := template.New("").Funcs(template.FuncMap{
		"seq":          seq,
		"contains":     strings.Contains,
		"emptyOrValue": emptyOrValue,
	}).ParseGlob("./../templates/*.html")

	if errTemplate != nil {
		fmt.Println("Erreur lors du chargement des templates :", errTemplate.Error())
		os.Exit(1)
	}

	http.Handle("/static/", http.StripPrefix("static/", http.FileServer(http.Dir("./../assets"))))

	game = Game{
		Grille:       initialiserGrille(),
		Joueur1:      Joueur{Pseudo: "", Couleur: "rouge"},
		Joueur2:      Joueur{Pseudo: "", Couleur: "jaune"},
		JoueurActuel: "rouge",
		Tour:         0,
		Message:      "Bienvenue ! Cliquez sur une colonne pour jouer.",
		Fin:          false,
		Gagnant:      Joueur{},
		Egalite:      false,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		listeTemplate.ExecuteTemplate(w, "index", nil)
	})

	http.HandleFunc("/templates/init", func(w http.ResponseWriter, r *http.Request) {
		listeTemplate.ExecuteTemplate(w, "init", nil)
	})

	http.HandleFunc("/templates/init/traitement", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			pseudo1 := r.FormValue("pseudo1")
			pseudo2 := r.FormValue("pseudo2")
			couleur1 := r.FormValue("couleur1")
			couleur2 := r.FormValue("couleur2")

			if couleur1 == couleur2 {
				data := struct {
					Error string
				}{
					Error: "Les deux joueurs ne peuvent pas avoir la même couleur",
				}
				listeTemplate.ExecuteTemplate(w, "init", data)
				return
			}

			if pseudo1 == pseudo2 {
				data := struct {
					Error string
				}{
					Error: "Les deux joueurs sont identiques",
				}
				listeTemplate.ExecuteTemplate(w, "init", data)
				return
			}

			game = Game{
				Grille:       initialiserGrille(),
				Joueur1:      Joueur{Pseudo: pseudo1, Couleur: couleur1},
				Joueur2:      Joueur{Pseudo: pseudo2, Couleur: couleur2},
				JoueurActuel: couleur1,
				Tour:         0,
				Message:      "La partie commence !",
				Fin:          false,
			}

			http.Redirect(w, r, "/templates/play", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/templates/play", func(w http.ResponseWriter, r *http.Request) {
		if game.Joueur1.Pseudo == "" || game.Joueur2.Pseudo == "" ||
			game.Joueur1.Pseudo == " " || game.Joueur2.Pseudo == " " {
			http.Redirect(w, r, "/templates/init", http.StatusSeeOther)
			return
		}

		if game.Fin {
			http.Redirect(w, r, "/templates/end", http.StatusSeeOther)
			return
		}

		listeTemplate.ExecuteTemplate(w, "play", game)
	})

	http.HandleFunc("/templates/play/traitement", func(w http.ResponseWriter, r *http.Request) {
		if game.Fin {
			http.Redirect(w, r, "/templates/end", http.StatusSeeOther)
			return
		}

		if r.Method == "POST" {
			colStr := r.FormValue("colonne")
			col, err := strconv.Atoi(colStr)
			if err != nil || col < 0 || col > 6 {
				game.Message = "Colonne invalide"
				http.Redirect(w, r, "/templates/play", http.StatusSeeOther)
				return
			}

			success := ajouterPion(&game.Grille, col, game.JoueurActuel)
			if success {

				vtype := verifierVictoireType(game.Grille, game.JoueurActuel)
				if vtype != "" {
					game.Fin = true
					game.VictoryType = vtype
					if game.JoueurActuel == game.Joueur1.Couleur {
						game.Gagnant = game.Joueur1
						enregistrerVictoire(game.Joueur1.Pseudo)
					} else {
						game.Gagnant = game.Joueur2
						enregistrerVictoire(game.Joueur2.Pseudo)
					}
					enregistrerHistorique(game.Joueur1.Pseudo, game.Joueur2.Pseudo, game.Gagnant.Pseudo)
					game.Message = fmt.Sprintf("Victoire de %s ! (%s)", game.Gagnant.Pseudo, vtype)

					http.Redirect(w, r, "/templates/end", http.StatusSeeOther)
					return
				} else if verifierEgalite(game.Grille) {
					game.Fin = true
					game.Egalite = true
					game.VictoryType = ""
					enregistrerHistorique(game.Joueur1.Pseudo, game.Joueur2.Pseudo, "")
					game.Message = "Égalité !"
					http.Redirect(w, r, "/templates/end", http.StatusSeeOther)
					return
				} else {
					game.JoueurActuel = switchJoueur(game.JoueurActuel)
					game.Tour++
					game.Message = "Coup joué avec succès"
				}
			} else {
				game.Message = "Colonne pleine"
			}
		}
		http.Redirect(w, r, "/templates/play", http.StatusSeeOther)
	})

	http.HandleFunc("/templates/end", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/templates/play", http.StatusSeeOther)
			return
		}
		if !game.Fin {
			http.Redirect(w, r, "/templates/play", http.StatusSeeOther)
			return
		}
		if err := listeTemplate.ExecuteTemplate(w, "end", game); err != nil {
			fmt.Println("Erreur template end:", err)
		}
	})

	fmt.Println("Serveur démarré sur : http://localhost:8080")
	http.HandleFunc("/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		data := ScoreboardData{
			Scores:     lireScores(),
			Historique: lireHistorique(),
		}

		listeTemplate.ExecuteTemplate(w, "scoreboard", data)
	})

	http.ListenAndServe(":8080", nil)

}

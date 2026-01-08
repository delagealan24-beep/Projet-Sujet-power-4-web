package main

import (
	"html/template"

	"log"

	"net/http"

	"strconv"
)

var templates *template.Template

func main() {
	var err error

	templates, err = template.ParseGlob("templates/*.gohtml")

	if err != nil {

		log.Fatal("Erreur chargement templates:", err)

	}
	http.HandleFunc("/health", healthHandler)

	http.HandleFunc("/games", gamesHandler)

	http.HandleFunc("/games/", singleGameHandler)

	http.HandleFunc("/play/", playHandler)

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Serveur démarré sur :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func playHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Path[len("/play/"):]

	id, err := strconv.Atoi(idStr)

	if err != nil {

		http.Error(w, "ID invalide", http.StatusBadRequest)

		return

	}

	gamesMutex.RLock()

	g, ok := games[id]

	gamesMutex.RUnlock()

	if !ok {

		http.Error(w, "Partie introuvable", http.StatusNotFound)

		return

	}

	g.mu.Lock()

	data := *g
	g.mu.Unlock()

	if err := templates.ExecuteTemplate(w, "game.gohtml", data); err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)

	}

}

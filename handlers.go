package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	pokemons, err := fetchPokemons()
	if err != nil {
		log.Printf("Erreur lors de la récupération des pokémons: %v", err)
		http.Error(w, "Impossible de charger les pokémons", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	data := struct {
		Title    string
		Pokemons []Pokemon
	}{
		Title:    "Pokédex - Tyradex",
		Pokemons: pokemons,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

func pokemonDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/pokemon/"):]
	idStr = strings.Split(idStr, "/")[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	pokemon, err := findPokemonByID(id)
	if err != nil {
		log.Printf("Erreur lors de la récupération du pokémon: %v", err)
		http.Error(w, "Impossible de charger le pokémon", http.StatusInternalServerError)
		return
	}

	if pokemon == nil {
		http.Error(w, "Pokémon introuvable", http.StatusNotFound)
		return
	}

	// Récupérer les notes depuis la base de données (READ)
	var notes []Note
	db.Where("pokemon_id = ?", id).Order("created_at desc").Find(&notes)

	funcMap := template.FuncMap{
		"divideBy": func(a int, b float64) float64 {
			return float64(a) / b * 100
		},
		"add": func(nums ...int) int {
			total := 0
			for _, n := range nums {
				total += n
			}
			return total
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006 à 15:04")
		},
	}

	tmpl := template.Must(template.New("detail.html").Funcs(funcMap).ParseFiles("templates/detail.html"))

	data := struct {
		Title   string
		Pokemon Pokemon
		Notes   []Note
	}{
		Title:   pokemon.Name.Fr + " - Pokédex",
		Pokemon: *pokemon,
		Notes:   notes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur template: %v", err)
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

// createNoteHandler crée une nouvelle note pour un Pokémon (CREATE)
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// URL : /pokemon/{id}/notes
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	pokemonID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur de formulaire", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, fmt.Sprintf("/pokemon/%d", pokemonID), http.StatusSeeOther)
		return
	}

	note := Note{
		PokemonID: uint(pokemonID),
		Content:   content,
	}

	// CREATE : insérer la note en base de données
	if result := db.Create(&note); result.Error != nil {
		log.Printf("Erreur création note: %v", result.Error)
		http.Error(w, "Impossible de créer la note", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pokemon/%d", pokemonID), http.StatusSeeOther)
}

// deleteNoteHandler supprime une note (DELETE)
func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// URL : /pokemon/{pokemonId}/notes/{noteId}/delete
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	pokemonID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "ID Pokémon invalide", http.StatusBadRequest)
		return
	}
	noteID, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		http.Error(w, "ID note invalide", http.StatusBadRequest)
		return
	}

	// DELETE : suppression de la note (soft delete via gorm.Model)
	if result := db.Delete(&Note{}, noteID); result.Error != nil {
		log.Printf("Erreur suppression note: %v", result.Error)
		http.Error(w, "Impossible de supprimer la note", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pokemon/%d", pokemonID), http.StatusSeeOther)
}

// editNoteHandler modifie le contenu d'une note (UPDATE)
func editNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// URL : /pokemon/{pokemonId}/notes/{noteId}/edit
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "URL invalide", http.StatusBadRequest)
		return
	}
	pokemonID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "ID Pokémon invalide", http.StatusBadRequest)
		return
	}
	noteID, err := strconv.ParseUint(parts[4], 10, 64)
	if err != nil {
		http.Error(w, "ID note invalide", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur de formulaire", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, fmt.Sprintf("/pokemon/%d", pokemonID), http.StatusSeeOther)
		return
	}

	// UPDATE : récupérer puis mettre à jour la note
	var note Note
	if result := db.First(&note, noteID); result.Error != nil {
		http.Error(w, "Note introuvable", http.StatusNotFound)
		return
	}

	note.Content = content
	if result := db.Save(&note); result.Error != nil {
		log.Printf("Erreur modification note: %v", result.Error)
		http.Error(w, "Impossible de modifier la note", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/pokemon/%d", pokemonID), http.StatusSeeOther)
}

// notesRouter dispatche les routes /pokemon/{id}/notes/*
func notesRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case strings.HasSuffix(path, "/notes") && r.Method == http.MethodPost:
		createNoteHandler(w, r)
	case strings.HasSuffix(path, "/delete") && r.Method == http.MethodPost:
		deleteNoteHandler(w, r)
	case strings.HasSuffix(path, "/edit") && r.Method == http.MethodPost:
		editNoteHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// pokemonRouter dispatche les routes /pokemon/*
func pokemonRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// Routes notes : /pokemon/{id}/notes ou /pokemon/{id}/notes/{noteId}/delete|edit
	if strings.Contains(path[len("/pokemon/"):], "/") {
		notesRouter(w, r)
		return
	}
	// Route détail : /pokemon/{id}
	pokemonDetailHandler(w, r)
}

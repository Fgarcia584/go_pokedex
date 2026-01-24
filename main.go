package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
    "os"
)

const (
	tyradexAPI = "https://tyradex.vercel.app/api/v1/pokemon"
)

// Pokemon représente un pokémon de l'API Tyradex
type Pokemon struct {
	PokedexID int    `json:"pokedex_id"`
	Name      Name   `json:"name"`
	Sprites   Sprites `json:"sprites"`
	Types     []Type `json:"types"`
}

type Name struct {
	Fr string `json:"fr"`
	En string `json:"en"`
	Jp string `json:"jp"`
}

type Sprites struct {
	Regular string `json:"regular"`
	Shiny   string `json:"shiny"`
}

type Type struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// Récupère la liste des pokémons depuis l'API Tyradex
func fetchPokemons() ([]Pokemon, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(tyradexAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pokemons []Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemons); err != nil {
		return nil, err
	}

	return pokemons, nil
}

// Handler pour la page d'accueil
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

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	// Routes
	http.HandleFunc("/", homeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
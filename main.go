package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"os"
)

const (
	tyradexAPI = "https://tyradex.vercel.app/api/v1/pokemon"
	port       = ":8080"
)

// Pokemon représente un pokémon de l'API Tyradex
type Pokemon struct {
	PokedexID    int         `json:"pokedex_id"`
	Generation   int         `json:"generation"`
	Category     string      `json:"category"`
	Name         Name        `json:"name"`
	Sprites      Sprites     `json:"sprites"`
	Types        []Type      `json:"types"`
	Talents      []Talent    `json:"talents"`
	Stats        Stats       `json:"stats"`
	Resistances  []Resistance `json:"resistances"`
	Evolution    *Evolution  `json:"evolution"`
	Height       string      `json:"height"`
	Weight       string      `json:"weight"`
	EggGroups    []string    `json:"egg_groups"`
	Sexe         *Sexe       `json:"sexe"`
	CatchRate    *int        `json:"catch_rate"`
	Level100     *int        `json:"level_100"`
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

type Talent struct {
	Name string `json:"name"`
	TC   bool   `json:"tc"`
}

type Stats struct {
	HP     int `json:"hp"`
	Attack int `json:"atk"`
	Def    int `json:"def"`
	SpAtk  int `json:"spe_atk"`
	SpDef  int `json:"spe_def"`
	Speed  int `json:"vit"`
}

type Resistance struct {
	Name       string  `json:"name"`
	Multiplier float64 `json:"multiplier"`
}

type Evolution struct {
	Pre  []EvolutionInfo `json:"pre"`
	Next []EvolutionInfo `json:"next"`
	Mega []EvolutionInfo `json:"mega"`
}

type EvolutionInfo struct {
	PokedexID int    `json:"pokedex_id"`
	Name      string `json:"name"`
	Condition string `json:"condition"`
}

type Sexe struct {
	Male   float64 `json:"male"`
	Female float64 `json:"female"`
}

var pokemonsCache []Pokemon
var cacheTime time.Time

// Récupère la liste des pokémons depuis l'API Tyradex avec cache
func fetchPokemons() ([]Pokemon, error) {
	// Cache de 5 minutes
	if time.Since(cacheTime) < 5*time.Minute && len(pokemonsCache) > 0 {
		return pokemonsCache, nil
	}

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

	pokemonsCache = pokemons
	cacheTime = time.Now()

	return pokemons, nil
}

// Trouve un Pokémon par son ID
func findPokemonByID(id int) (*Pokemon, error) {
	pokemons, err := fetchPokemons()
	if err != nil {
		return nil, err
	}

	for _, p := range pokemons {
		if p.PokedexID == id {
			return &p, nil
		}
	}

	return nil, nil
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

// Handler pour la page de détails d'un Pokémon
func pokemonDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'URL
	idStr := r.URL.Path[len("/pokemon/"):]
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

	// Fonctions helper pour le template
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
	}

	tmpl := template.Must(template.New("detail.html").Funcs(funcMap).ParseFiles("templates/detail.html"))
	
	data := struct {
		Title   string
		Pokemon Pokemon
	}{
		Title:   pokemon.Name.Fr + " - Pokédex",
		Pokemon: *pokemon,
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
	http.HandleFunc("/pokemon/", pokemonDetailHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
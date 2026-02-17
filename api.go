package main

import (
	"encoding/json"
	"net/http"
	"time"
)

var pokemonsCache []Pokemon
var cacheTime time.Time

// fetchPokemons récupère la liste des Pokémon depuis l'API Tyradex avec un cache de 5 minutes
func fetchPokemons() ([]Pokemon, error) {
	if time.Since(cacheTime) < 5*time.Minute && len(pokemonsCache) > 0 {
		return pokemonsCache, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}

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

// findPokemonByID recherche un Pokémon par son identifiant dans le cache
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

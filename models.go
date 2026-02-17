package main

import "gorm.io/gorm"

var tyradexAPI string

// -------------------------
// Modèle GORM
// -------------------------

// Note représente un commentaire laissé sur un Pokémon
type Note struct {
	gorm.Model
	PokemonID uint   `gorm:"not null;index"`
	Content   string `gorm:"not null"`
}

// -------------------------
// Modèles API Tyradex
// -------------------------

type Pokemon struct {
	PokedexID   int          `json:"pokedex_id"`
	Generation  int          `json:"generation"`
	Category    string       `json:"category"`
	Name        Name         `json:"name"`
	Sprites     Sprites      `json:"sprites"`
	Types       []Type       `json:"types"`
	Talents     []Talent     `json:"talents"`
	Stats       Stats        `json:"stats"`
	Resistances []Resistance `json:"resistances"`
	Evolution   *Evolution   `json:"evolution"`
	Height      string       `json:"height"`
	Weight      string       `json:"weight"`
	EggGroups   []string     `json:"egg_groups"`
	Sexe        *Sexe        `json:"sexe"`
	CatchRate   *int         `json:"catch_rate"`
	Level100    *int         `json:"level_100"`
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

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run in development
go run main.go

# Build for deployment
go build -o bin/go_pokedex main.go

# Format code
go fmt ./...

# Vet code
go vet ./...

# Run tests (no tests currently exist)
go test ./...
```

The server listens on port `8080` by default, configurable via the `PORT` environment variable.

## Architecture

Single-file Go web application (`main.go`) with zero external dependencies. Uses only the standard library.

**Routes:**
- `GET /` → `homeHandler` — fetches all Pokemon, renders `templates/index.html`
- `GET /pokemon/{id}` → `pokemonDetailHandler` — finds Pokemon by ID from cache, renders `templates/detail.html`
- `GET /static/*` → file server for `static/` directory

**Data fetching:**
Pokemon data comes from the Tyradex API (`https://tyradex.app/api/v1/pokemon`). Results are stored in a package-level `pokemonsCache` variable with a 5-minute TTL managed by `lastFetchTime`. The `fetchPokemons()` function handles cache checks and API calls. Individual Pokemon lookup uses `findPokemonByID()` which searches the cache by `PokedexID`.

**Templates** in `templates/` use Go's `html/template`. The index page includes embedded JavaScript for keyboard navigation (arrow keys, Enter) and click-based Pokemon selection with live display updates in the left panel.

**Deployment:** The `Procfile` targets `./bin/go_pokedex` for Scalingo/Heroku-style platforms.

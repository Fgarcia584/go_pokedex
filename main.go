package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	tyradexAPI = os.Getenv("TYRADEX_API_URL")
	if tyradexAPI == "" {
		log.Fatal("La variable d'environnement TYRADEX_API_URL est requise")
	}

	initDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/pokemon/", pokemonRouter)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Serveur démarré sur le port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

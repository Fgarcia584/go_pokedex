package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("La variable d'environnement DATABASE_URL est requise")
	}

	// DEBUG: affiche la connection string complète (retire ça après)
    log.Printf("CONNECTION STRING: %s", dsn)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Impossible de se connecter à la base de données : %v", err)
	}

	// AutoMigrate crée la table si elle n'existe pas
	if err := db.AutoMigrate(&Note{}); err != nil {
		log.Fatalf("Erreur lors de la migration : %v", err)
	}

	log.Println("Base de données connectée et migrée avec succès")
}

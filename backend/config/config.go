package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Printf("Greška prilikom učitavanja .env fajla: %v", err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
		log.Println("SERVER_PORT nije pronađen, koristi se podrazumevani port 8080")
	}

	return &Config{
		ServerPort: port,
	}
}

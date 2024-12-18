package main

import (
	ws "backend-ZI/websocket"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	http.HandleFunc("/ws", ws.WsHandler)
	log.Println("WebSocket server pokrenut na :8080")
	serverRun()
}

func serverRun() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Greška prilikom učitavanja .env fajla: %v", err)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
		fmt.Println("SERVER_PORT nije pronađen, koristi se podrazumevani port 8080")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Zdravo! Server radi na portu %s", port)
	})

	address := fmt.Sprintf(":%s", port)
	fmt.Printf("Server se pokreće na http://localhost%s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("Greška prilikom pokretanja servera: %v", err)
	}
}

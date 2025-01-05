package main

import (
	"backend-ZI/config"
	"backend-ZI/services"
	"backend-ZI/websocket"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	controlChannel := make(chan string)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, controlChannel)
	})
	http.HandleFunc("/upload", services.EnableCORS(services.UploadHandler))

	http.HandleFunc("/control", services.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		services.CommandHandler(w, r, controlChannel)
	}))

	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running at http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

package main

import (
	"backend-ZI/config"
	"backend-ZI/filewatcher"
	"backend-ZI/services"
	"backend-ZI/websocket"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	controlChannel := make(chan string)
	events := make(chan []string)
	dirToWatch := "/Users/dusanpavlovic016/Books/Target"

	log.Println("Pokrenuta gorutina za file Target")
	go filewatcher.WatchDir(dirToWatch, events, controlChannel)

	controlChannelX := make(chan string)
	eventsX := make(chan []string)
	dirToWatchX := "/Users/dusanpavlovic016/Books/X"

	log.Println("Pokrenuta gorutina za file X")
	go filewatcher.WatchDir(dirToWatchX, eventsX, controlChannelX)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, controlChannel, events, controlChannelX, eventsX)
	})

	http.HandleFunc("/upload", services.EnableCORS(services.UploadHandler))
	http.HandleFunc("/uploadandencode", services.EnableCORS(services.UploadandencodeHandler))
	http.HandleFunc("/encodetype", services.EnableCORS(services.EncodeTypeHandler))

	http.HandleFunc("/control", services.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		services.CommandHandler(w, r, controlChannel, controlChannelX)
	}))

	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running at http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

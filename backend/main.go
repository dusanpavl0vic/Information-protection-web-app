package main

import (
	"backend-ZI/config"
	"backend-ZI/filewatcher"
	"backend-ZI/services"
	"backend-ZI/websocket"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	cfg := config.LoadConfig()

	controlChannel := make(chan string)
	events := make(chan []string)
	dirToWatch := "/Users/dusanpavlovic016/Books/Target"

	var watchOnce sync.Once
	watchOnce.Do(func() {
		log.Println("Pokrenuta gorutina za file Target")
		go filewatcher.WatchDir(dirToWatch, events, controlChannel)
	})
	controlChannelX := make(chan string)
	eventsX := make(chan []string)
	dirToWatchX := "/Users/dusanpavlovic016/Books/X"

	var watchOnce2 sync.Once
	watchOnce2.Do(func() {
		log.Println("Pokrenuta gorutina za file X")
		go filewatcher.WatchDir2(dirToWatchX, eventsX, controlChannelX)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(w, r, controlChannel, events, controlChannelX, eventsX)
	})

	http.HandleFunc("/upload", services.EnableCORS(services.UploadHandler))
	http.HandleFunc("/uploadandencode", services.EnableCORS(services.UploadandencodeHandler))
	http.HandleFunc("/encodetype", services.EnableCORS(services.CipherTypeHandler))
	http.HandleFunc("/file-list-x-action", services.EnableCORS(services.DecodeFileHandler))

	http.HandleFunc("/control", services.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		services.CommandHandler(w, r, controlChannel, controlChannelX)
	}))

	address := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running at http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

package websocket

import (
	"backend-ZI/filewatcher"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var activeConnections int32
var watchOnce sync.Once

func HandleWebSocket(w http.ResponseWriter, r *http.Request, controlChannel chan string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	atomic.AddInt32(&activeConnections, 1)
	log.Printf("WebSocket connection established. Active connections: %d", activeConnections)

	defer func() {
		atomic.AddInt32(&activeConnections, -1)
		conn.Close()
		log.Printf("WebSocket connection closed. Active connections: %d", activeConnections)
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Initial deadline
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Refresh on pong
		return nil
	})

	events := make(chan []string)
	dirToWatch := "/Users/dusanpavlovic016/Books/Target"

	watchOnce.Do(func() {
		log.Println("Pokrenuta gorutina za file watcher")
		go filewatcher.WatchDir(dirToWatch, events, controlChannel)
	})

	for files := range events {
		if err := conn.WriteJSON(files); err != nil {
			log.Println("Error sending data to WebSocket client:", err)
			return
		}
	}
}

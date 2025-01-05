package websocket

import (
	"backend-ZI/filewatcher"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (for production, implement specific checks)
	},
}
var activeConnections int32
var watchOnce sync.Once

// HandleWebSocket sets up a WebSocket connection and sends data when receiving a file list.
func HandleWebSocket(w http.ResponseWriter, r *http.Request, controlChannel chan string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	atomic.AddInt32(&activeConnections, 1)
	log.Printf("WebSocket connection established. Active connections: %d", activeConnections)

	defer func() {
		conn.Close()
		atomic.AddInt32(&activeConnections, -1)
		log.Printf("WebSocket connection closed. Active connections: %d", activeConnections)
	}()
	events := make(chan []string)
	log.Println("WebSocket connection established")

	dirToWatch := "/Users/dusanpavlovic016/Books/Target"

	watchOnce.Do(func() {
		go func() {
			log.Println("Watching directory: ./Target")
			filewatcher.WatchDir(dirToWatch, events, controlChannel)
		}()
	})

	for files := range events {
		log.Printf("Sending file list to client: %v\n", files)
		err := conn.WriteJSON(files)
		if err != nil {
			log.Println("Error sending data to WebSocket client:", err)
			return
		}
	}

}

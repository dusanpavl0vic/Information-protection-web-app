package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var activeConnections = make(map[*websocket.Conn]bool)

func HandleWebSocket(w http.ResponseWriter, r *http.Request, controlChannel chan string, events chan []string, controlChannelX chan string, eventsX chan []string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Dodaj novu konekciju u mapu aktivnih konekcija
	activeConnections[conn] = true
	log.Printf("WebSocket connection established. Active connections: %d", len(activeConnections))

	defer func() {
		// Ukloni konekciju kada se zatvori
		delete(activeConnections, conn)
		log.Printf("WebSocket connection closed. Active connections: %d", len(activeConnections))
	}()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message or connection closed:", err)
				// Kada dođe do greške, zatvori konekciju i izađi iz gorutine
				conn.Close()
				return
			}
		}
	}()

	for {
		select {
		case files := <-events:
			message := WebSocketMessage{Type: "eventFiles", Data: files} // Kreiraj poruku sa tipom "eventFiles"
			for c := range activeConnections {
				mutex.Lock()
				err := c.WriteJSON(message)
				mutex.Unlock()
				if err != nil {
					log.Printf("Ugasen WebSocket klijent. Zatvaram kone")
					delete(activeConnections, c)
					c.Close()
				}
			}
		case filesX := <-eventsX:
			message := WebSocketMessage{Type: "eventFilesX", Data: filesX} // Kreiraj poruku sa tipom "eventFilesX"
			for c := range activeConnections {
				mutex.Lock()
				err := c.WriteJSON(message)
				mutex.Unlock()
				if err != nil {
					log.Printf("Ugasen WebSocket klijent. Zatvaram kone")
					delete(activeConnections, c)
					c.Close()
				}
			}
		}
	}
}

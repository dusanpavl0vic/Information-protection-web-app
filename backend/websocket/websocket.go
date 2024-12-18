package websocket

import (
	"backend-ZI/filewatcher"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func WsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Greška pri upgrade-u WebSocket-a:", err)
		return
	}
	defer conn.Close()

	updateChan := make(chan []string)

	dirToWatch := "./files"

	go filewatcher.FileWatcher(dirToWatch, updateChan)

	for files := range updateChan {
		log.Println("Slanje fajlova klijentu:", files)
		if err := conn.WriteJSON(files); err != nil {
			log.Println("Greška pri slanju podataka klijentu:", err)
			break // Zatvara konekciju ako se desi greška
		}
	}
	log.Println("WebSocket konekcija je zatvorena sa strane servera.")

}

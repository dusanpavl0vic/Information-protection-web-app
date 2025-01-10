package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CommandRequest struct {
	Command string `json:"command"`
}

func CommandHandler(w http.ResponseWriter, r *http.Request, controlChannel chan string, controlChannelX chan string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	switch req.Command {
	case "start":
		controlChannel <- "start"
		controlChannelX <- "start"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"success": true, "message": "File watcher started."}`)
	case "stop":
		controlChannel <- "stop"
		controlChannelX <- "stop"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"success": true, "message": "File watcher stopped."}`)
	default:
		http.Error(w, `{"success": false, "message": "Invalid command."}`, http.StatusBadRequest)
	}
}

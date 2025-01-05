package filewatcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchDir(dir string, events chan []string, controlChannel chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return // Umesto log.Fatalf, koristi povratak
	}

	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		log.Printf("Failed to watch directory %s: %v", dir, err)
		return // Nastavlja rad bez prekida celog servera
	}

	//log.Printf("Watching directory: %s\n", dir)

	active := true // Početno stanje: aktivan
	for {
		select {
		case event := <-watcher.Events:
			if active && event.Op&(fsnotify.Create) != 0 {
				log.Printf("File created in directory: %s\n", dir)
				files, err := listFiles(dir)
				if err != nil {
					log.Printf("Error listing files in directory %s: %v", dir, err)
					continue
				}
				events <- files
			}
		case err := <-watcher.Errors:
			log.Printf("File watcher error: %v", err)
		case command := <-controlChannel:
			switch command {
			case "stop":
				log.Println("Stopping file watcher...")
				active = false
			case "start":
				log.Println("Starting file watcher...")
				active = true
				files, err := listFiles(dir)
				if err != nil {
					log.Printf("Error listing files in directory %s: %v", dir, err)
					continue
				}
				events <- files
			}
		}
	}
}

func listFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

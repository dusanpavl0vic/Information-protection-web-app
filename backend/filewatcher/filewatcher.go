package filewatcher

import (
	"backend-ZI/services"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchDir(dir string, events chan []string, controlChannel chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}

	defer watcher.Close()

	err = watcher.Add(dir)
	if err != nil {
		log.Printf("Failed to watch directory %s: %v", dir, err)
		return
	}

	active := false
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
				go func(filePath string) {
					fileName := filepath.Base(filePath)
					fileData, err1 := readFile(filePath)
					if err1 != nil {
						log.Printf("Error reading file %s: %v", filePath, err1)
						return
					}
					err2 := services.EncodeFile(fileData, fileName)
					if err2 != nil {
						log.Printf("Error encoding file %s: %v", fileName, err2)
					}
				}(event.Name)
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
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	return files, err
}

func readFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	return data, nil
}

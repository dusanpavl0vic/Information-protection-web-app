package filewatcher

import (
	"backend-ZI/coders"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type CipherType int

const (
	RailFence CipherType = iota // Automatski dobija vrednost 0
	XXTEA                       // Automatski dobija vrednost 1
)

func (c CipherType) String() string {
	switch c {
	case RailFence:
		return "RailFence"
	case XXTEA:
		return "XXTEA"
	default:
		return "UnknownCipher"
	}
}

var (
	depth  int        = 3
	cipher CipherType = 0
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

	active := true
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
				go func(fileName string) {
					err2 := encodeFile(fileName)
					if err2 != nil {
						log.Printf("Error encoding file %s: %v", fileName, err2) // Popravka da koristi err2 umesto err
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
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func encodeFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var encrypted string

	switch cipher {
	case RailFence:
		encrypted = coders.EncryptRailFence(string(data), depth)
		fmt.Println("Cipher is Railfence")
	case XXTEA:
		fmt.Println("Cipher is XXTEA")
	}

	modifiedPath := strings.Replace(file, "/Target/", "/X/", 1)

	dir, file := filepath.Split(modifiedPath)
	ext := filepath.Ext(file)
	baseName := strings.TrimSuffix(file, ext)

	newFileName := fmt.Sprintf("%s_%s%s", baseName, cipher.String(), ext)
	newFilePath := filepath.Join(dir, newFileName)

	err2 := os.WriteFile(newFilePath, []byte(encrypted), 0644)
	if err2 != nil {
		return fmt.Errorf("failed to write encrypted file: %v", err)
	}

	fmt.Println("File saved successfully:", newFilePath)
	return nil
}

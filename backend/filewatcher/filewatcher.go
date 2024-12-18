package filewatcher

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func listFiles(dir string) ([]string, error) {

	var files []string
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}
	// indeks, pojedinacni element
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name()) // dodajemo u listu
		}
	}

	return files, nil

}

func FileWatcher(dirToWatch string, updateChan chan<- []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Greška prilikom kreiranja Watchera: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(dirToWatch)
	if err != nil {
		log.Fatalf("Greška prilikom dodavanja direktorijuma u Watcher: %v", err)
	}

	// pocetna lista fajlova salje websocketu
	sendFileList(dirToWatch, updateChan)

	// proveravamo file
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
				sendFileList(dirToWatch, updateChan)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Greška u Watcheru: %v", err)
		}
	}
}

// salje listu fajlova prko kanala
// komunikacija izmedju gorutina
func sendFileList(dir string, updateChan chan<- []string) {
	files, err := listFiles(dir)
	if err != nil {
		log.Printf("Greška pri listanju fajlova: %v", err)
		return
	}
	updateChan <- files
}

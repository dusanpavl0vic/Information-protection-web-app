package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func saveFile(file multipart.File, fileHeader *multipart.FileHeader, targetDir string) (string, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("nije moguće kreirati direktorijum: %v", err)
		}
	}

	filePath := filepath.Join(targetDir, fileHeader.Filename)

	destFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("nije moguće kreirati fajl na disku: %v", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("nije moguće sačuvati fajl: %v", err)
	}

	return filePath, nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metoda nije podržana", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("nije moguće obraditi fajl: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("mene niko ne zove")

	targetDir := "/Users/dusanpavlovic016/Books/Target"

	filePath, err := saveFile(file, fileHeader, targetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Fajl je uspešno sačuvan na: %s", filePath)
}

func UploadandencodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metoda nije podržana", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	log.Println(fileHeader.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("nije moguće obraditi fajl: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("nije moguće pročitati fajl: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("update i ekoduj samo se ja javljam babbabababa")

	EncodeFile(fileData, fileHeader.Filename)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "File uploaded and encoded successfully!"}`))
}

package services

import (
	"backend-ZI/coders"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	mu     sync.Mutex
)

func EncodeFile(fileData []byte, filename string) error {

	var encrypted string

	switch cipher {
	case RailFence:
		encrypted = coders.EncryptRailFence(string(fileData), depth)
		fmt.Println("Cipher is Railfence")
	case XXTEA:
		fmt.Println("Cipher is XXTEA")
	}

	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	newFileName := fmt.Sprintf("%s_%s%s", baseName, cipher.String(), ext)
	newFilePath := filepath.Join("/Users/dusanpavlovic016/Books/X", newFileName)

	err2 := os.WriteFile(newFilePath, []byte(encrypted), 0644)
	if err2 != nil {
		return fmt.Errorf("failed to write encrypted file: %v", err2)
	}

	fmt.Println("File saved successfully:", newFilePath)
	return nil
}

func EncodeTypeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metoda nije podržana", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Cipher CipherType `json:"cipher"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Neispravan zahtev: %v", err), http.StatusBadRequest)
		return
	}

	mu.Lock()
	updateCipher(&cipher, requestData.Cipher)
	mu.Unlock()

	log.Printf("Odabrani algoritam: %s", cipher.String())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Selected cipher: %s"}`, cipher.String())))
}

// Funkcija za ažuriranje vrednosti pomoću pokazivača
func updateCipher(currentCipher *CipherType, newCipher CipherType) {
	*currentCipher = newCipher
}

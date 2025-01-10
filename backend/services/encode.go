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

type DecodeRequest struct {
	FileName string `json:"fileName"`
}

type DecodeResponse struct {
	DecodedFileName string `json:"decodedFileName"`
	DecodedContent  string `json:"decodedContent"`
}

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

	log.Println("cao ovde sam prosoo 1")
	fmt.Println("File saved successfully:", newFilePath)
	return nil
}

func CipherTypeHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Println("updatovan chiper 1")

	mu.Unlock()

	log.Printf("Odabrani algoritam: %s", cipher.String())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Selected cipher: %s"}`, cipher.String())))
}

func updateCipher(currentCipher *CipherType, newCipher CipherType) {
	*currentCipher = newCipher
}

func DecodeFile(fileName string) ([]byte, string, error) {
	// Putanja do fajla
	filePath := filepath.Join("/Users/dusanpavlovic016/Books/X", fileName)

	// Pročitaj sadržaj fajla
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}

	var decoded string

	// Dekodiranje na osnovu algoritma
	switch cipher {
	case RailFence:
		decoded = coders.DecryptRailFence(string(fileData), depth)
		fmt.Println("Cipher is RailFence", decoded)
	case XXTEA:
		fmt.Println("Cipher is XXTEA")
		// decodedData, err := coders.DecryptXXTEA(fileData, xxteaKey)
		// if err != nil {
		// 	return nil, "", fmt.Errorf("failed to decode with XXTEA: %v", err)
		// }
		// decoded = string(decodedData)
	default:
		return nil, "", fmt.Errorf("unsupported cipher")
	}

	// Dodaj "decoded" u naziv fajla
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)
	decodedFileName := fmt.Sprintf("%s_decoded%s", baseName, ext)

	fmt.Println("Decoded file name:", decodedFileName)
	return []byte(decoded), decodedFileName, nil
}

func DecodeFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req DecodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	decodedData, decodedFileName, err := DecodeFile(req.FileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode file: %v", err), http.StatusInternalServerError)
		return
	}

	resp := DecodeResponse{
		DecodedFileName: decodedFileName,
		DecodedContent:  string(decodedData),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		return
	}
}

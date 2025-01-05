package services

import "net/http"

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Dodaj CORS zaglavlja
		w.Header().Set("Access-Control-Allow-Origin", "*") // Omogućava sve izvore
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Ako je preflight OPTIONS zahtev, odmah odgovori sa 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r) // Pozovi originalni handler
	}
}

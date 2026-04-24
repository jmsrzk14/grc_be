package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// respondJSON mengirimkan response HTTP dalam format JSON.
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	log.NewHelper(log.DefaultLogger).Infof("Response sent: %d bytes", len(response))
}

// respondError mengirimkan response HTTP error dalam format JSON.
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// parseUUIDFromRequest mengekstrak variabel UUID dari path menggunakan Gorilla Mux vars.
func parseUUIDFromRequest(r *http.Request, key string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	val := vars[key]
	return uuid.Parse(val)
}

// parseUUID memparsing string menjadi uuid.UUID.
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

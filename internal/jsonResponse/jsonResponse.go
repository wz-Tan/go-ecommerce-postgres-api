package jsonResponse

import (
	"encoding/json"
	"net/http"
)

// General Function for Response Handling
func Write(w http.ResponseWriter, status int, data any) {
	// Headers for Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // Locks Header

	// Turn Data into Bytes
	json.NewEncoder(w).Encode(data)
}

package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// Success writes an success response to a request in json format.
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	JSON(w, statusCode, data)
}

// Error writes an error response in json format to the request.
func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

// JSON writes the response to a request in json format.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if error := json.NewEncoder(w).Encode(data); error != nil {
			log.Fatal(error)
		}
	}
}

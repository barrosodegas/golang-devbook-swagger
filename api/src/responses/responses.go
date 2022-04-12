package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

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

// Error writes an error response in json format to the request.
func Error(w http.ResponseWriter, statusCode int, er error) {
	JSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: er.Error(),
	})
}

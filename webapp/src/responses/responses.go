package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != http.StatusNoContent {
		if error := json.NewEncoder(w).Encode(data); error != nil {
			log.Fatal(error)
		}
	}
}

func PrepareAndReturnError(w http.ResponseWriter, r *http.Response) {
	var error ResponseError
	json.NewDecoder(r.Body).Decode(&error)
	JSON(w, r.StatusCode, error)
}

package middleware

import (
	"api/src/authentication"
	"api/src/responses"
	"fmt"
	"net/http"
)

// Logger logs data about the request made by the user.
// returns the next function to be executed.
func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Method=%s RequestURI=%s Host=%s\n", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}

// Authenticate authenticates the user's request before processing the request.
// returns the next function to be executed.
func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if error := authentication.ValidateToken(r); error != nil {
			responses.Error(w, http.StatusUnauthorized, error)
			return
		}
		nextFunction(w, r)
	}
}

package middleware

import (
	"fmt"
	"log"
	"net/http"
	"webapp/src/cookies"
)

func Logger(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method=%s RequestURI=%s Host=%s\n", r.Method, r.RequestURI, r.Host)
		nextFunction(w, r)
	}
}

func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, error := cookies.Read(r); error != nil {
			fmt.Printf("Authetication error: %s\n", error)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		nextFunction(w, r)
	}
}

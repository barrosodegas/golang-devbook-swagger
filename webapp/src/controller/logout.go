package controller

import (
	"net/http"
	"webapp/src/cookies"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookies.ClearCookies(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

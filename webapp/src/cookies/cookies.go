package cookies

import (
	"net/http"
	"time"
	"webapp/src/config"

	"github.com/gorilla/securecookie"
)

var s *securecookie.SecureCookie

func ConfigureSecureCookie() {
	s = securecookie.New(config.HashKey, config.BlockKey)
}

func SaveCookies(w http.ResponseWriter, ID, token string) error {
	devBookdata := map[string]string{
		"id":    ID,
		"token": token,
	}

	dataEncoded, error := s.Encode("devbook-data", devBookdata)
	if error != nil {
		return error
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "devbook-data",
		Value:    dataEncoded,
		Path:     "/",
		HttpOnly: true,
	})

	return nil
}

func Read(r *http.Request) (map[string]string, error) {
	cookie, error := r.Cookie("devbook-data")
	if error != nil {
		return nil, error
	}

	cookieValues := make(map[string]string)
	error = s.Decode("devbook-data", cookie.Value, &cookieValues)
	if error != nil {
		return nil, error
	}

	return cookieValues, nil
}

func ClearCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "devbook-data",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}

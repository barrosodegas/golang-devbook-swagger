package authentication

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateToken creates a jwt token.
// return a jwt token or an error if an error occurs during the process.
func CreateToken(userId uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString(config.SecretKey)
}

// ValidateToken validates the given token.
// return an error if the token is invalid or return nil if it is valid.
func ValidateToken(r *http.Request) error {
	textToken := extractToken(r)

	token, error := jwt.Parse(textToken, getVerificationKey)
	if error != nil {
		fmt.Printf("Invalid token with error: %s\n", error)
		return errors.New("Invalid token!")
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Invalid token!")
}

// ExtractUserIdOfToken extract the user id from the token
// returns the userid or an error.
func ExtractUserIdOfToken(r *http.Request) (uint64, error) {
	textToken := extractToken(r)

	token, error := jwt.Parse(textToken, getVerificationKey)
	if error != nil {
		fmt.Printf("Invalid token with error: %s\n", error)
		return 0, errors.New("Invalid token!")
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, error := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["userId"]), 10, 64)
		if error != nil {
			return 0, error
		}
		return userId, nil
	}

	return 0, errors.New("Invalid token!")
}

// extractToken extracts the request token.
// returns the request token.
func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	words := strings.Split(strings.TrimSpace(token), " ")
	if len(words) == 2 {
		return words[1]
	}
	return ""
}

// getVerificationKey retrieves the verification key.
// returns the verification key or an error if the signing method is different than expected.
func getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}

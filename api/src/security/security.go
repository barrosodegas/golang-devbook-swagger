package security

import "golang.org/x/crypto/bcrypt"

// GenerateHash generates the hash for the user's password.
// returns the encrypted password or an error.
func GenerateHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword checks if the entered password is valid.
// returns an error if the password is invalid.
func VerifyPassword(passwordHash, passwordText string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordText))
}

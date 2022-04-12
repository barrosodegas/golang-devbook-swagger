package model

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// User represents a user.
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// Prepare validates and clears the fields of an user.
// returns an error if the user is invalid.
func (u *User) Prepare(step string) error {
	if error := u.validate(step); error != nil {
		return error
	}

	u.clearFields(step)

	if step == "create" {
		passwordHash, error := security.GenerateHash(u.Password)
		if error != nil {
			return error
		}
		u.Password = string(passwordHash)
	}
	return nil
}

// validate validate an user.
// returns an error if the user is invalid.
func (u *User) validate(step string) error {
	if u.Name == "" {
		return errors.New("Name is required!")
	}

	if u.Nick == "" {
		return errors.New("Nick is required!")
	}

	if u.Email == "" {
		return errors.New("Email is required!")
	}

	if error := checkmail.ValidateFormat(u.Email); error != nil {
		return errors.New("Email is invalid!")
	}

	if step == "create" && u.Password == "" {
		return errors.New("Password is required!")
	}

	return nil
}

// clearFields clears a publication fields.
func (u *User) clearFields(step string) {
	u.Name = strings.TrimSpace(u.Name)
	u.Nick = strings.TrimSpace(u.Nick)
	u.Email = strings.TrimSpace(u.Email)
}

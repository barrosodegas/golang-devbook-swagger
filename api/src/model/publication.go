package model

import (
	"errors"
	"strings"
	"time"
)

// Publication represents a publication.
type Publication struct {
	ID         uint64    `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content,omitempty"`
	AuthorId   uint64    `json:"authorId,omitempty"`
	AuthorNick string    `json:"authorNick,omitempty"`
	Likes      uint64    `json:"likes"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

// Prepare validates and clears the fields of a publication.
// returns an error if the publication is invalid.
func (p *Publication) Prepare() error {
	if error := p.validate(); error != nil {
		return error
	}

	p.clearFields()

	return nil
}

// validate validate a publication.
// returns an error if the publication is invalid.
func (p *Publication) validate() error {
	if p.Title == "" {
		return errors.New("Title is required!")
	}

	if p.Content == "" {
		return errors.New("Content is required!")
	}

	return nil
}

// clearFields clears a publication fields
func (p *Publication) clearFields() {
	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
}

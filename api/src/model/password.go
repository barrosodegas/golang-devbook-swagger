package model

// Password represents the user's passwords that are in the process of being updated.
type Password struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

package models

type Author struct {
	ID        int    `json:",omitempty"`
	FirstName string `json:",omitempty"`
	LastName  string `json:",omitempty"`
	Dob       string `json:",omitempty"`
	PenName   string `json:",omitempty"`
}
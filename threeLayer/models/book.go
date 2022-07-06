package models

type Book struct {
	ID              int    `json:",omitempty"`
	Title           string `json:",omitempty"`
	Author          Author `json:",omitempty"`
	Publication     string `json:",omitempty"`
	PublicationDate string `json:",omitempty"`
}

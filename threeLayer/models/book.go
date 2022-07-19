package models

type Book struct {
	ID              int    `json:",omitempty"`
	Title           string `json:",omitempty"`
	Author          Author `json:",omitempty"`
	Publication     string `json:",omitempty"`
	PublicationDate string `json:",omitempty"`
}

type contextKey string

const (
	Title         contextKey = "title"
	IncludeAuthor contextKey = "includeAuthor"
)

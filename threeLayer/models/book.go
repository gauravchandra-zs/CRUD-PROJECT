package models

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          Author `json:"author,omitempty"`
	Publication     string `json:"publication"`
	PublicationDate string `json:"publication_date"`
}

type contextKey string

const (
	Title         contextKey = "title"
	IncludeAuthor contextKey = "includeAuthor"
)

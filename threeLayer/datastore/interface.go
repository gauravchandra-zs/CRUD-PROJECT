package datastore

import (
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type Book interface {
	GetAllBooks(title string) ([]models.Book, error)
	GetBookByID(id int) (models.Book, error)
	PostBook(book *models.Book) (int, error)
	DeleteBook(id int) (int, error)
	PutBook(book *models.Book) (models.Book, error)
	DeleteBookByAuthorID(id int) error
}

type Author interface {
	PostAuthor(author models.Author) (int, error)
	DeleteAuthor(id int) (int, error)
	PutAuthor(id int, author models.Author) (models.Author, error)
	GetAuthorByID(id int) (models.Author, error)
}

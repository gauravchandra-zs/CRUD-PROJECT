package service

import (
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type Book interface {
	GetAllBooks(params map[string]string) ([]models.Book, error)
	GetBookByID(id int) (models.Book, error)
	PostBook(book *models.Book) (int, error)
	DeleteBook(id int) (int, error)
	PutBook(book *models.Book) (models.Book, error)
}

type Author interface {
	PostAuthor(author models.Author) (int, error)
	DeleteAuthor(id int) (int, error)
	PutAuthor(id int, author models.Author) (models.Author, error)
}

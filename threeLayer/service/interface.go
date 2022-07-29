package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type Book interface {
	GetAllBooks(ctx *gofr.Context) ([]models.Book, error)
	GetBookByID(ctx *gofr.Context, id int) (models.Book, error)
	PostBook(ctx *gofr.Context, book *models.Book) (int, error)
	DeleteBook(ctx *gofr.Context, id int) (int, error)
	PutBook(ctx *gofr.Context, id int, book *models.Book) (models.Book, error)
}

type Author interface {
	PostAuthor(ctx *gofr.Context, author models.Author) (int, error)
	DeleteAuthor(ctx *gofr.Context, id int) (int, error)
	PutAuthor(ctx *gofr.Context, id int, author models.Author) (models.Author, error)
}

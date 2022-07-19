package datastore

import (
	"context"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type Book interface {
	GetAllBooks(ctx *gofr.Context) ([]models.Book, error)
	GetAllBooksByTitle(ctx *gofr.Context, title string) ([]models.Book, error)
	GetBookByID(ctx *gofr.Context, id int) (models.Book, error)
	PostBook(ctx *gofr.Context, book *models.Book) (int, error)
	DeleteBook(ctx context.Context, id int) (int, error)
	PutBook(ctx *gofr.Context, id int, book *models.Book) (models.Book, error)
	DeleteBookByAuthorID(ctx context.Context, id int) error
	CheckBook(ctx context.Context, book *models.Book) bool
	CheckBookBid(ctx context.Context, id int) bool
}

type Author interface {
	PostAuthor(ctx context.Context, author models.Author) (int, error)
	DeleteAuthor(ctx context.Context, id int) (int, error)
	PutAuthor(ctx *gofr.Context, id int, author models.Author) (models.Author, error)
	GetAuthorByID(ctx *gofr.Context, id int) (models.Author, error)
	CheckAuthorByID(ctx context.Context, id int) bool
	CheckAuthor(ctx context.Context, author models.Author) bool
}

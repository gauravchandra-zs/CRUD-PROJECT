package datastore

import (
	"context"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type Book interface {
	GetAllBooks(ctx context.Context, title string) ([]models.Book, error)
	GetAllBooksByTitle(ctx context.Context, title string) ([]models.Book, error)
	GetBookByID(ctx context.Context, id int) (models.Book, error)
	PostBook(ctx context.Context, book *models.Book) (int, error)
	DeleteBook(ctx context.Context, id int) (int, error)
	PutBook(ctx context.Context, id int, book *models.Book) (models.Book, error)
	DeleteBookByAuthorID(ctx context.Context, id int) error
	CheckBook(ctx context.Context, book *models.Book) bool
	CheckBookBid(ctx context.Context, id int) bool
}

type Author interface {
	PostAuthor(ctx context.Context, author models.Author) (int, error)
	DeleteAuthor(ctx context.Context, id int) (int, error)
	PutAuthor(ctx context.Context, id int, author models.Author) (models.Author, error)
	GetAuthorByID(ctx context.Context, id int) (models.Author, error)
	CheckAuthorByID(ctx context.Context, id int) bool
	CheckAuthor(ctx context.Context, author models.Author) bool
}

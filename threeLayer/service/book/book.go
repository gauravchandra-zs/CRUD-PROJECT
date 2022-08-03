package servicebook

import (
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"MyProject/CRUD-PROJECT/threeLayer/datastore"
	"MyProject/CRUD-PROJECT/threeLayer/models"
)

type ServiceBook struct {
	bookStore   datastore.Book
	authorStore datastore.Author
}

func New(bookStore datastore.Book, author datastore.Author) ServiceBook {
	return ServiceBook{bookStore, author}
}

// GetAllBooks call GetAllBooks and GetAllBooksByTitle of store layer to get all details
func (s ServiceBook) GetAllBooks(ctx *gofr.Context) ([]models.Book, error) {
	title := ctx.Param(string(models.Title))
	includeAuthor := ctx.Param(string(models.IncludeAuthor))

	var output []models.Book

	var err error

	if title == "" {
		output, err = s.bookStore.GetAllBooks(ctx)
	} else {
		output, err = s.bookStore.GetAllBooksByTitle(ctx, title)
	}

	if err != nil {
		return output, err
	}

	for i := 0; i < len(output); i++ {
		var author models.Author

		if includeAuthor == "true" {
			author, err = s.authorStore.GetAuthorByID(ctx, output[i].Author.ID)
			if err != nil {
				return output, err
			}
		}

		output[i].Author = author
	}

	return output, nil
}

// GetBookByID check book exist with given id and call GetBookByID on store layer to get Book detail
func (s ServiceBook) GetBookByID(ctx *gofr.Context, id int) (models.Book, error) {
	var output models.Book

	var err error

	output, err = s.bookStore.GetBookByID(ctx, id)
	if err != nil || output.Author.ID <= 0 {
		return output, err
	}

	output.Author, err = s.authorStore.GetAuthorByID(ctx, output.Author.ID)
	if err != nil {
		return output, err
	}

	return output, nil
}

// PostBook validate book and author detail and call PostBook on store layer to post book
func (s ServiceBook) PostBook(ctx *gofr.Context, book *models.Book) (int, error) {
	if s.bookStore.CheckBook(ctx, book) || !s.authorStore.CheckAuthorByID(ctx, book.Author.ID) {
		return 0, &errors.Response{StatusCode: http.StatusConflict, Reason: "entity already exists", Code: "Entity Already Exists"}
	}

	id, err := s.bookStore.PostBook(ctx, book)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteBook check book exist with given id and call DeleteBook on store layer to delete book
func (s ServiceBook) DeleteBook(ctx *gofr.Context, id int) (int, error) {
	if !s.bookStore.CheckBookBid(ctx, id) {
		return 0, errors.EntityNotFound{}
	}

	return s.bookStore.DeleteBook(ctx, id)
}

// PutBook validate author and book detail and call PutAuthor on store layer to update book detail
func (s ServiceBook) PutBook(ctx *gofr.Context, id int, book *models.Book) (models.Book, error) {
	var output models.Book

	if !s.authorStore.CheckAuthorByID(ctx, book.Author.ID) || !s.bookStore.CheckBookBid(ctx, id) {
		return models.Book{}, errors.EntityNotFound{}
	}

	author, err := s.authorStore.PutAuthor(ctx, book.Author.ID, book.Author)
	if err != nil {
		return models.Book{}, err
	}

	book.Author = author

	output, err = s.bookStore.PutBook(ctx, id, book)
	if err != nil {
		return models.Book{}, err
	}

	return output, nil
}

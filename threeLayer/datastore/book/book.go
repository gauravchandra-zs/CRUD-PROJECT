package datastorebook

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"MyProject/CRUD-PROJECT/threeLayer/drivers"
	"MyProject/CRUD-PROJECT/threeLayer/models"
)

type BookStore struct {
}

func New() BookStore {
	return BookStore{}
}

// GetAllBooksByTitle return all book detail from Book table with given title
func (b BookStore) GetAllBooksByTitle(ctx *gofr.Context, title string) ([]models.Book, error) {
	var books []models.Book

	result, err := ctx.DB().QueryContext(ctx, drivers.SelectFromBookByTitle, title)
	if err != nil {
		return []models.Book{}, errors.EntityNotFound{}
	}

	check := true
	for result.Next() {
		check = false
		var book models.Book

		err = result.Scan(&book.ID, &book.Title, &book.Publication, &book.PublicationDate, &book.Author.ID)
		if err != nil {
			return []models.Book{}, err
		}

		books = append(books, book)
	}

	if check {
		return nil, errors.EntityNotFound{}
	}

	return books, nil
}

// GetAllBooks  return all books from Book table
func (b BookStore) GetAllBooks(ctx *gofr.Context) ([]models.Book, error) {
	var output []models.Book

	result, err := ctx.DB().QueryContext(ctx, drivers.SelectFromBook)
	if err != nil {
		return []models.Book{}, errors.EntityNotFound{}
	}

	for result.Next() {
		var book models.Book

		err = result.Scan(&book.ID, &book.Title, &book.Publication, &book.PublicationDate, &book.Author.ID)
		if err != nil {
			return []models.Book{}, err
		}

		output = append(output, book)
	}

	return output, nil
}

// GetBookByID return Book with given id from book table
func (b BookStore) GetBookByID(ctx *gofr.Context, id int) (models.Book, error) {
	var book models.Book

	result, err := ctx.DB().QueryContext(ctx, drivers.SelectFromBookByID, id)
	if err != nil {
		return book, err
	}

	if result.Next() {
		b := models.Book{}

		err = result.Scan(&b.ID, &b.Title, &b.Publication, &b.PublicationDate, &b.Author.ID)
		if err != nil {
			return book, err
		}

		return b, nil
	}

	return book, errors.EntityNotFound{}
}

// PostBook post book with given detail if not exist
func (b BookStore) PostBook(ctx *gofr.Context, book *models.Book) (int, error) {
	rs, err := ctx.DB().DB.ExecContext(ctx, drivers.InsertIntoBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil {
		return 0, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

// DeleteBook delete book detail with given id
func (b BookStore) DeleteBook(ctx *gofr.Context, id int) (int, error) {
	result, err := ctx.DB().ExecContext(ctx, drivers.DeleteBookQuery, id)
	if err != nil {
		return 0, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(r), nil
}

// PutBook update book detail if book exist with given id and detail
func (b BookStore) PutBook(ctx *gofr.Context, id int, book *models.Book) (models.Book, error) {
	_, err := ctx.DB().ExecContext(ctx, drivers.UpdateBook, book.Title, book.Publication, book.PublicationDate, id)
	if err != nil {
		return models.Book{}, errors.Error("error in updating book")
	}

	return *book, nil
}

// DeleteBookByAuthorID delete book with given authorID if exist
func (b BookStore) DeleteBookByAuthorID(ctx *gofr.Context, id int) error {
	_, err := ctx.DB().ExecContext(ctx, drivers.DeleteBookByAuthorID, id)
	if err != nil {
		return err
	}

	return nil
}

// CheckBook check book exist or not with given detail and return bool value
func (b BookStore) CheckBook(ctx *gofr.Context, book *models.Book) bool {
	result, err := ctx.DB().QueryContext(ctx, drivers.CheckBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil || !result.Next() {
		return false
	}

	return true
}

// CheckBookBid check book exist  or not with given id and return bool value
func (b BookStore) CheckBookBid(ctx *gofr.Context, id int) bool {
	res, err := ctx.DB().QueryContext(ctx, drivers.CheckBookBYID, id)
	if err != nil || !res.Next() {
		return false
	}

	return true
}

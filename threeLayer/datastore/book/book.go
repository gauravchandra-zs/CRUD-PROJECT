package datastorebook

import (
	"context"
	"database/sql"
	"errors"

	"Projects/GoLang-Interns-2022/threeLayer/drivers"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type BookStore struct {
	db *sql.DB
}

func New(db *sql.DB) BookStore {
	return BookStore{db: db}
}

func (b BookStore) GetAllBooksByTitle(ctx context.Context, title string) ([]models.Book, error) {
	var output []models.Book

	result, err := b.db.QueryContext(ctx, drivers.SelectFromBookByTitle, title)
	if err != nil {
		return []models.Book{}, nil
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
func (b BookStore) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	var output []models.Book

	result, err := b.db.QueryContext(ctx, drivers.SelectFromBook)
	if err != nil {
		return []models.Book{}, nil
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

func (b BookStore) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	var book models.Book

	result, err := b.db.QueryContext(ctx, drivers.SelectFromBookByID, id)
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

	return book, errors.New("book not exists")
}

func (b BookStore) PostBook(ctx context.Context, book *models.Book) (int, error) {
	rs, err := b.db.ExecContext(ctx, drivers.InsertIntoBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil {
		return 0, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (b BookStore) DeleteBook(ctx context.Context, id int) (int, error) {
	result, err := b.db.ExecContext(ctx, drivers.DeleteBookQuery, id)
	if err != nil {
		return 0, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(r), nil
}

func (b BookStore) PutBook(ctx context.Context, id int, book *models.Book) (models.Book, error) {
	_, err := b.db.ExecContext(ctx, drivers.UpdateBook, book.Title, book.Publication, book.PublicationDate, id)
	if err != nil {
		return models.Book{}, errors.New("error in updating book")
	}

	return *book, nil
}

func (b BookStore) DeleteBookByAuthorID(ctx context.Context, id int) error {
	_, err := b.db.ExecContext(ctx, drivers.DeleteBookByAuthorID, id)
	if err != nil {
		return err
	}

	return nil
}

func (b BookStore) CheckBook(ctx context.Context, book *models.Book) bool {
	result, err := b.db.QueryContext(ctx, drivers.CheckBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil || !result.Next() {
		return false
	}

	return true
}

func (b BookStore) CheckBookBid(ctx context.Context, id int) bool {
	res, err := b.db.QueryContext(ctx, drivers.CheckBookBYID, id)
	if err != nil || !res.Next() {
		return false
	}

	return true
}

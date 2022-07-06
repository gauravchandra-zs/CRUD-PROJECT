package datastorebook

import (
	"database/sql"
	"errors"

	"Projects/GoLang-Interns-2022/threeLayer/driver"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type BookStore struct {
	db *sql.DB
}

func New(db *sql.DB) BookStore {
	return BookStore{db: db}
}

func (b BookStore) GetAllBooks(title string) ([]models.Book, error) {
	var output []models.Book

	var result *sql.Rows

	var err error

	if title == "" {
		result, err = b.db.Query(driver.SelectFromBook)
	} else {
		result, err = b.db.Query(driver.SelectFromBookByTitle, title)
	}

	if err != nil {
		return output, nil
	}

	for result.Next() {
		var book models.Book

		err = result.Scan(&book.ID, &book.Title, &book.Publication, &book.PublicationDate, &book.Author.ID)
		if err != nil {
			return output, err
		}

		output = append(output, book)
	}

	return output, nil
}

func (b BookStore) GetBookByID(id int) (models.Book, error) {
	var book models.Book

	result, err := b.db.Query(driver.SelectFromBookByID, id)
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

func (b BookStore) PostBook(book *models.Book) (int, error) {
	result, err := b.db.Query(driver.CheckBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil || result.Next() {
		return 0, err
	}

	rs, err := b.db.Exec(driver.InsertIntoBook, book.Title, book.Publication, book.PublicationDate, book.Author.ID)
	if err != nil {
		return 0, err
	}

	id, err := rs.LastInsertId()

	return int(id), err
}

func (b BookStore) DeleteBook(id int) (int, error) {
	var row int64

	res, err := b.db.Query(driver.CheckBookBYID, id)
	if !res.Next() || err != nil {
		return int(row), err
	}

	result, err := b.db.Exec(driver.DeleteBookQuery, id)
	if err != nil {
		return int(row), err
	}

	row, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(row), nil
}

func (b BookStore) PutBook(book *models.Book) (models.Book, error) {
	var output models.Book

	if !b.checkBookBid(book.ID) {
		return output, errors.New("author or book not exits")
	}

	_, err := b.db.Exec(driver.UpdateBook, book.Title, book.Publication, book.PublicationDate, book.ID)
	if err != nil {
		return output, errors.New("error in updating book")
	}

	return *book, nil
}

func (b BookStore) checkBookBid(id int) bool {
	res, err := b.db.Query(driver.CheckBookBYID, id)
	if !res.Next() || err != nil {
		return false
	}

	return true
}

func (b BookStore) GetAuthorByID(id int) (models.Author, error) {
	var output models.Author

	ResAuthor, err := b.db.Query(driver.SelectAuthorByID, id)
	if err != nil {
		return output, err
	}

	if ResAuthor.Next() {
		err := ResAuthor.Scan(&output.ID, &output.FirstName, &output.LastName, &output.Dob, &output.PenName)
		if err != nil {
			return models.Author{}, err
		}
	}

	return output, nil
}

func (b BookStore) DeleteBookByAuthorID(id int) error {
	_, err := b.db.Exec("DELETE FROM Book WHERE AuthorID=?", id)
	if err != nil {
		return err
	}

	return nil
}

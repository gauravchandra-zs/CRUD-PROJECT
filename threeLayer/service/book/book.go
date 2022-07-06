package servicebook

import (
	"errors"
	"log"
	"strings"

	"Projects/GoLang-Interns-2022/threeLayer/datastore"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type ServiceBook struct {
	bookStore   datastore.Book
	authorStore datastore.Author
}

func New(bookStore datastore.Book, author datastore.Author) ServiceBook {
	return ServiceBook{bookStore, author}
}

func (s ServiceBook) GetAllBooks(params map[string]string) ([]models.Book, error) {
	var output []models.Book

	title := params["title"]
	authorInclude := params["authorInclude"]

	output, err := s.bookStore.GetAllBooks(title)
	if err != nil {
		return output, err
	}

	for i := 0; i < len(output); i++ {
		var author models.Author

		if authorInclude == "true" {
			author, err = s.authorStore.GetAuthorByID(output[i].Author.ID)
			if err != nil {
				return output, err
			}
		}

		output[i].Author = author
	}

	return output, nil
}

func (s ServiceBook) GetBookByID(id int) (models.Book, error) {
	var output models.Book

	var err error

	if id <= 0 {
		return output, errors.New("invalid id")
	}

	output, err = s.bookStore.GetBookByID(id)
	if err != nil || output.Author.ID <= 0 {
		return output, err
	}

	output.Author, err = s.authorStore.GetAuthorByID(output.Author.ID)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (s ServiceBook) PostBook(book *models.Book) (int, error) {
	if !validateBook(book) || !validateAuthor(book.Author) {
		return 0, errors.New("invalid book or author")
	}

	insertedID, err := s.authorStore.PostAuthor(book.Author)
	if err != nil || insertedID <= 0 {
		log.Print(err)
	}

	book.Author.ID = insertedID

	insertedID, err = s.bookStore.PostBook(book)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (s ServiceBook) DeleteBook(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	deletedID, err := s.bookStore.DeleteBook(id)
	if err != nil || deletedID <= 0 {
		return deletedID, errors.New("unsuccessful deletion")
	}

	return deletedID, nil
}

func (s ServiceBook) PutBook(book *models.Book) (models.Book, error) {
	var output models.Book

	if !validateBook(book) || !validateAuthor(book.Author) {
		return output, errors.New("invalid book or author")
	}

	_, err := s.authorStore.PutAuthor(book.Author.ID, book.Author)
	if err != nil {
		return models.Book{}, err
	}

	output, err = s.bookStore.PutBook(book)
	if err != nil {
		return output, err
	}

	return output, nil
}

func validateBook(b *models.Book) bool {
	slc := strings.Split(b.PublicationDate, "-")
	sz := 3

	switch {
	case b.Publication != "Scholastic" && b.Publication != "Penguin" && b.Publication != "Arihanth":
		return false
	case len(slc) < sz:
		return false
	case slc[2] >= "2022" || slc[2] < "1880":
		return false
	case b.Title == "":
		return false
	default:
		return true
	}
}

func validateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

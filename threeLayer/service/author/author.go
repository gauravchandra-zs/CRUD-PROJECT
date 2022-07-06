package serviceauthor

import (
	"errors"

	"Projects/GoLang-Interns-2022/threeLayer/datastore"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type ServiceAuthor struct {
	authorStore datastore.Author
	bookStore   datastore.Book
}

func New(author datastore.Author, book datastore.Book) ServiceAuthor {
	return ServiceAuthor{author, book}
}

func (s ServiceAuthor) PostAuthor(author models.Author) (int, error) {
	if !ValidateAuthor(author) {
		return 0, errors.New("invalid author detail")
	}

	insertedID, err := s.authorStore.PostAuthor(author)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (s ServiceAuthor) PutAuthor(id int, author models.Author) (models.Author, error) {
	var output models.Author

	if !ValidateAuthor(author) {
		return output, errors.New("invalid author detail")
	}

	output, err := s.authorStore.PutAuthor(id, author)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (s ServiceAuthor) DeleteAuthor(id int) (int, error) {
	err := s.bookStore.DeleteBookByAuthorID(id)
	if err != nil {
		return 0, err
	}

	id, err = s.authorStore.DeleteAuthor(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ValidateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

package serviceauthor

import (
	"context"
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

// PostAuthor validate author detail and call for postAuthor to store layer
func (s ServiceAuthor) PostAuthor(ctx context.Context, author models.Author) (int, error) {
	if s.authorStore.CheckAuthor(ctx, author) {
		return 0, errors.New("author exist")
	}

	insertedID, err := s.authorStore.PostAuthor(ctx, author)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

// PutAuthor  check author exist or not and call PutAuthor on store layer to update author
func (s ServiceAuthor) PutAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	if !s.authorStore.CheckAuthorByID(ctx, id) {
		return models.Author{}, errors.New("author not exist")
	}

	output, err := s.authorStore.PutAuthor(ctx, id, author)
	if err != nil {
		return output, err
	}

	return output, nil
}

// DeleteAuthor check author exist or not and call DeleteAuthor on store layer to delete author
func (s ServiceAuthor) DeleteAuthor(ctx context.Context, id int) (int, error) {
	if !s.authorStore.CheckAuthorByID(ctx, id) {
		return 0, errors.New("author not exist")
	}

	err := s.bookStore.DeleteBookByAuthorID(ctx, id)
	if err != nil {
		return 0, err
	}

	id, err = s.authorStore.DeleteAuthor(ctx, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

package servicebook

import (
	"context"
	"errors"
	"log"

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

func (s ServiceBook) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	title, ok := ctx.Value("title").(string)
	if !ok {
		log.Print("title is not of type string")
	}

	includeAuthor, ok := ctx.Value("includeAuthor").(string)
	if !ok {
		log.Print("includeAuthor is not of type string")
	}

	output, err := s.bookStore.GetAllBooks(ctx, title)
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

func (s ServiceBook) GetBookByID(ctx context.Context, id int) (models.Book, error) {
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

func (s ServiceBook) PostBook(ctx context.Context, book *models.Book) (int, error) {
	if s.bookStore.CheckBook(ctx, book) || !s.authorStore.CheckAuthorByID(ctx, book.Author.ID) {
		return 0, errors.New("book exist already")
	}

	id, err := s.bookStore.PostBook(ctx, book)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	return id, nil
}

func (s ServiceBook) DeleteBook(ctx context.Context, id int) (int, error) {
	if !s.bookStore.CheckBookBid(ctx, id) {
		return 0, errors.New("book not exist")
	}

	rowEffected, err := s.bookStore.DeleteBook(ctx, id)
	if err != nil {
		return rowEffected, errors.New("unsuccessful deletion")
	}

	return rowEffected, nil
}

func (s ServiceBook) PutBook(ctx context.Context, id int, book *models.Book) (models.Book, error) {
	var output models.Book

	if !s.authorStore.CheckAuthorByID(ctx, book.Author.ID) {
		return models.Book{}, errors.New("author not present")
	}

	author, err := s.authorStore.PutAuthor(ctx, book.Author.ID, book.Author)
	if err != nil {
		return models.Book{}, err
	}

	book.Author = author

	if !s.bookStore.CheckBookBid(ctx, id) {
		return output, err
	}

	output, err = s.bookStore.PutBook(ctx, id, book)
	if err != nil {
		return output, err
	}

	return output, nil
}

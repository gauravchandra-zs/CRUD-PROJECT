package serviceauthor

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type mockAuthorStore struct {
}

func (m mockAuthorStore) CheckAuthorByID(ctx context.Context, id int) bool {
	return true
}

func (m mockAuthorStore) CheckAuthor(ctx context.Context, author models.Author) bool {
	return true
}

func (m mockAuthorStore) PostAuthor(ctx context.Context, author models.Author) (int, error) {
	if author.FirstName == "gaurav" && author.LastName == "chandra" {
		return 0, errors.New("author already exist")
	}

	return 1, nil
}

func (m mockAuthorStore) DeleteAuthor(ctx context.Context, id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	if id == 10 {
		return 0, errors.New("author not exist")
	}

	return 1, nil
}

func (m mockAuthorStore) PutAuthor(ctx context.Context, id int, author models.Author) (models.Author, error) {
	if id <= 0 || id >= 100 {
		return models.Author{}, errors.New("invalid author id")
	}

	return author, nil
}

func (m mockAuthorStore) GetAuthorByID(ctx context.Context, id int) (models.Author, error) {
	return models.Author{}, nil
}

type mockBookStore struct {
}

func (m mockBookStore) CheckBook(ctx context.Context, book *models.Book) bool {
	return true
}

func (m mockBookStore) CheckBookBid(ctx context.Context, id int) bool {
	return true
}

func (m mockBookStore) GetAllBooks(ctx context.Context, title string) ([]models.Book, error) {
	return nil, nil
}

func (m mockBookStore) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	return models.Book{}, nil
}

func (m mockBookStore) PostBook(ctx context.Context, book *models.Book) (int, error) {
	return 0, nil
}

func (m mockBookStore) DeleteBook(ctx context.Context, id int) (int, error) {
	return 0, nil
}

func (m mockBookStore) PutBook(ctx context.Context, id int, book *models.Book) (models.Book, error) {
	return models.Book{}, nil
}

func (m mockBookStore) DeleteBookByAuthorID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	if id == 5 {
		return errors.New("book not exist")
	}

	return nil
}

func TestServiceAuthor_PostAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		body       models.Author
		expectedID int
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 1,
		},
		{
			"invalid case", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0,
		},
		{
			"invalid case author exist", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0,
		},
	}
	for _, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})
		ctx := context.Background()
		id, _ := a.PostAuthor(ctx, v.body)

		if !reflect.DeepEqual(id, v.expectedID) {
			t.Errorf("Expected %v\tGot %v", v.expectedID, id)
		}
	}
}

func TestServiceAuthor_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		body           models.Author
		expectedOutput models.Author
	}{
		{
			"valid case", 1, models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
		},

		{
			"invalid id", 0, models.Author{
				ID:        0,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
		},
		{
			"author not exist", 100, models.Author{
				ID:        100,
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
		},
	}
	for _, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})
		ctx := context.Background()
		output, _ := a.PutAuthor(ctx, v.id, v.body)

		if !reflect.DeepEqual(v.expectedOutput, output) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestServiceAuthor_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		deletedID int
	}{
		{
			"valid case", 1, 1,
		},
		{
			"invalid case", -1, 0,
		},
		{
			"valid case", 5, 0,
		},
		{
			"valid case", 10, 0,
		},
	}
	for _, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})
		ctx := context.Background()
		id, _ := a.DeleteAuthor(ctx, v.id)

		if id != v.deletedID {
			t.Errorf("test case fail")
		}
	}
}

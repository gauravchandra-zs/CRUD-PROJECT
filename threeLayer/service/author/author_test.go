package serviceauthor

import (
	"errors"
	"log"
	"reflect"
	"testing"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type mockAuthorStore struct {
}

func (m mockAuthorStore) PostAuthor(author models.Author) (int, error) {
	if author.FirstName == "gaurav" && author.LastName == "chandra" {
		return 0, errors.New("author already exist")
	}

	return 1, nil
}

func (m mockAuthorStore) DeleteAuthor(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	if id == 10 {
		return 0, errors.New("author not exist")
	}

	return 1, nil
}

func (m mockAuthorStore) PutAuthor(id int, author models.Author) (models.Author, error) {
	if id <= 0 || id >= 100 {
		return models.Author{}, errors.New("invalid author id")
	}

	return author, nil
}

func (m mockAuthorStore) GetAuthorByID(id int) (models.Author, error) {
	return models.Author{}, nil
}

type mockBookStore struct {
}

func (m mockBookStore) GetAllBooks(title string) ([]models.Book, error) {
	return nil, nil
}

func (m mockBookStore) GetBookByID(id int) (models.Book, error) {
	return models.Book{}, nil
}

func (m mockBookStore) PostBook(book *models.Book) (int, error) {
	return 0, nil
}

func (m mockBookStore) DeleteBook(id int) (int, error) {
	return 0, nil
}

func (m mockBookStore) PutBook(book *models.Book) (models.Book, error) {
	return models.Book{}, nil
}

func (m mockBookStore) DeleteBookByAuthorID(id int) error {
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

		id, err := a.PostAuthor(v.body)
		if err != nil {
			log.Print(err)
		}

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

		output, err := a.PutAuthor(v.id, v.body)
		if err != nil {
			log.Print(err)
		}

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
			"author exist but book not", 5, 0,
		},
		{
			"valid author but not exist", 10, 0,
		},
	}
	for _, v := range testcases {
		a := New(mockAuthorStore{}, mockBookStore{})

		id, err := a.DeleteAuthor(v.id)
		if err != nil {
			log.Print(err)
		}

		if id != v.deletedID {
			t.Errorf("test case fail")
		}
	}
}

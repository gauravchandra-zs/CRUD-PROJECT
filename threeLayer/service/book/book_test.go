package servicebook

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"Projects/GoLang-Interns-2022/threeLayer/models"
)

type mockBookstore struct{}

func (m mockBookstore) CheckBook(ctx context.Context, book *models.Book) bool {
	//TODO implement me
	panic("implement me")
}

func (m mockBookstore) CheckBookBid(ctx context.Context, id int) bool {
	//TODO implement me
	panic("implement me")
}

func (m mockBookstore) GetAllBooks(ctx context.Context, title string) ([]models.Book, error) {
	if title == "" {
		return []models.Book{
			{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
			{
				ID:    2,
				Title: "RD",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		}, nil
	} else if title == "RD" {
		return []models.Book{
			{
				ID:    2,
				Title: "RD",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		}, nil
	} else if title == "xyz" {
		return []models.Book{}, errors.New("invalid title")
	}

	return []models.Book{}, nil
}

func (m mockBookstore) GetBookByID(id int) (models.Book, error) {
	if id <= 0 {
		return models.Book{}, errors.New("invalid id")
	}

	if id == 1 {
		return models.Book{
			ID:    1,
			Title: "RD sharma",
			Author: models.Author{
				ID: 1,
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		}, nil
	}

	if id == 2 {
		return models.Book{
			ID:    2,
			Title: "RD sharma",
			Author: models.Author{
				ID: 10,
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		}, nil
	}

	return models.Book{}, errors.New("book not exist")
}

func (m mockBookstore) PostBook(book *models.Book) (int, error) {
	if book.Title == "RD" && book.Publication == "Penguin" && book.PublicationDate == "12-08-2011" {
		return 0, errors.New("book exist")
	}

	return 2, nil
}

func (m mockBookstore) DeleteBook(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	return 1, nil
}

func (m mockBookstore) PutBook(book *models.Book) (models.Book, error) {
	if book.ID <= 0 {
		return models.Book{}, errors.New("invalid Book id")
	}

	return *book, nil
}

func (m mockBookstore) DeleteBookByAuthorID(id int) error {
	return nil
}

type mockAuthorStore struct{}

func (m mockAuthorStore) CheckAuthorByID(ctx context.Context, id int) bool {
	//TODO implement me
	panic("implement me")
}

func (m mockAuthorStore) CheckAuthor(ctx context.Context, author models.Author) bool {

}

func (m mockAuthorStore) PostAuthor(author models.Author) (int, error) {
	if author.FirstName == "gaurav" && author.LastName == "chaudhari" {
		return 0, errors.New("author already exist")
	}

	return 1, nil
}

func (m mockAuthorStore) DeleteAuthor(id int) (int, error) {
	return 0, nil
}

func (m mockAuthorStore) PutAuthor(id int, author models.Author) (models.Author, error) {
	if id <= 0 {
		return models.Author{}, errors.New("invalid author id")
	}

	return author, nil
}

func (m mockAuthorStore) GetAuthorByID(ctx context.Context, id int) (models.Author, error) {
	if id <= 0 {
		return models.Author{}, errors.New("invalid id")
	}

	if id == 1 {
		return models.Author{
			ID:        1,
			FirstName: "gaurav",
			LastName:  "chandra",
			Dob:       "18-07-2001",
			PenName:   "GCC",
		}, nil
	}

	return models.Author{}, errors.New("author not exist")
}

func TestGetAllBooks(t *testing.T) {
	testcases := []struct {
		desc           string
		params         map[string]string
		expectedOutput []models.Book
	}{
		{
			"valid case", map[string]string{
			"title": "", "includeAuthor": "true"}, []models.Book{
			{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID:        1,
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
			{
				ID:    2,
				Title: "RD",
				Author: models.Author{
					ID:        1,
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		},
		},
		{
			"valid case", map[string]string{
			"title": "", "includeAuthor": ""}, []models.Book{
			{
				ID:              1,
				Title:           "RD sharma",
				Author:          models.Author{},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
			{
				ID:              2,
				Title:           "RD",
				Author:          models.Author{},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		},
		},
		{
			"invalid case", map[string]string{
			"title": "xyz", "authorInclude": ""},
			[]models.Book{},
		},
	}
	for _, v := range testcases {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "title", v.params["title"])
		ctx = context.WithValue(ctx, "includeAuthor", v.params["includeAuthor"])
		a := New(mockBookstore{}, mockAuthorStore{})
		output, _ := a.GetAllBooks(ctx)

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestGetBookByID(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput models.Book
	}{
		{
			"valid case", 1, models.Book{
			ID:    1,
			Title: "RD sharma",
			Author: models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		},
		},
		{"invalid case", -1, models.Book{}},
		{"invalid case", 100, models.Book{}},
		{
			"invalid case", 2, models.Book{
			ID:              2,
			Title:           "RD sharma",
			Author:          models.Author{},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		},
		},
	}
	for _, v := range testcases {
		a := New(mockBookstore{}, mockAuthorStore{})
		ctx := context.Background()
		output, _ := a.GetBookByID(ctx, v.id)

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc           string
		body           models.Book
		lastInsertedID int
	}{
		{
			"valid case", models.Book{
			Title: "RD sharma",
			Author: models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		}, 2,
		},

		{
			"invalid book already exists", models.Book{
			Title: "RD",
			Author: models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "Penguin",
			PublicationDate: "12-08-2011",
		}, 0,
		},
		{
			"invalid case", models.Book{
			Title: "",
			Author: models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "NCERT",
			PublicationDate: "12-08-2011",
		}, 0,
		},
	}
	for _, v := range testcases {
		a := New(mockBookstore{}, mockAuthorStore{})
		ctx := context.Background()

		id, _ := a.PostBook(ctx, &v.body)

		if !reflect.DeepEqual(id, v.lastInsertedID) {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		body      models.Book
		expOutput models.Book
	}{
		{
			"valid case", models.Book{
			ID:    1,
			Title: "RD sharma",
			Author: models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		}, models.Book{
			ID:    1,
			Title: "RD sharma",
			Author: models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		},
		},
		{
			"invalid case", models.Book{
			ID:    -1,
			Title: "",
			Author: models.Author{
				ID:        0,
				FirstName: "gaurav",
				LastName:  "",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
			Publication:     "NCERT",
			PublicationDate: "12-08-2011",
		}, models.Book{},
		},
	}
	for _, v := range testcases {
		a := New(mockBookstore{}, mockAuthorStore{})
		ctx := context.Background()
		output, _ := a.PutBook(ctx, &v.body)

		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		deletedID int
	}{
		{
			"valid case", 1, 1,
		},
		{
			"valid case", -1, 0,
		},
	}
	for _, v := range testcases {
		a := New(mockBookstore{}, mockAuthorStore{})
		ctx := context.Background()

		deletedID, _ := a.DeleteBook(ctx, v.id)

		if !reflect.DeepEqual(v.deletedID, deletedID) {
			t.Errorf("Expected %v\tGot %v", v.deletedID, deletedID)
		}
	}
}

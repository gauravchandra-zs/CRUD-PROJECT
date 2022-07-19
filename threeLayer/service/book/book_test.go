package servicebook

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"Projects/GoLang-Interns-2022/threeLayer/datastore"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

func TestGetAllBooks(t *testing.T) {
	testcases := []struct {
		desc           string
		params         map[string]any
		expectedOutput []models.Book
		err            error
	}{
		{
			"valid case", map[string]any{
				"title": "", "includeAuthor": "true"}, []models.Book{
				{
					ID:          1,
					Title:       "RD sharma",
					Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
					Publication: "Arihanth", PublicationDate: "12-08-2011",
				},
				{
					ID:          2,
					Title:       "RD",
					Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
					Publication: "Arihanth", PublicationDate: "12-08-2011",
				},
			}, errors.New("author not exist"),
		},
		{
			"invalid case", map[string]any{
				"title": 1, "includeAuthor": true}, []models.Book{
				{
					ID:          1,
					Title:       "RD sharma",
					Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
					Publication: "Arihanth", PublicationDate: "12-08-2011",
				},
				{
					ID:          2,
					Title:       "RD",
					Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
					Publication: "Arihanth", PublicationDate: "12-08-2011",
				},
			}, nil,
		},
		{
			"invalid case error from GetAllBook", map[string]any{
				"title": "xyz", "includeAuthor": "true"}, []models.Book{},
			errors.New("book not exists"),
		},
	}

	var author = models.Author{
		ID:        1,
		FirstName: "gaurav",
		LastName:  "chandra",
		Dob:       "18-07-2001",
		PenName:   "GCC",
	}

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		ctx := context.Background()
		ctx = context.WithValue(ctx, models.Title, v.params["title"])
		ctx = context.WithValue(ctx, models.IncludeAuthor, v.params["includeAuthor"])

		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)

		mockBookStore.EXPECT().GetAllBooks(ctx).Return(v.expectedOutput, nil).AnyTimes()
		mockBookStore.EXPECT().GetAllBooksByTitle(ctx, v.params["title"]).Return(v.expectedOutput, v.err).AnyTimes()
		mockAuthorStore.EXPECT().GetAuthorByID(ctx, author.ID).Return(author, v.err).AnyTimes()

		a := New(mockBookStore, mockAuthorStore)
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
		errBook        error
		errAuthor      error
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
			}, nil, nil,
		},
		{
			"invalid case error from getBook", 100, models.Book{},
			errors.New("book not exist"), nil,
		},
		{
			"invalid case error from getAuthor", 2, models.Book{ID: 1, Author: models.Author{
				ID: 1,
			}}, nil, errors.New("author not exist"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		ctx := context.Background()

		mockBookStore.EXPECT().GetBookByID(ctx, v.id).Return(v.expectedOutput, v.errBook)
		mockAuthorStore.EXPECT().GetAuthorByID(ctx, v.expectedOutput.Author.ID).Return(v.expectedOutput.Author, v.errAuthor).AnyTimes()
		a := New(mockBookStore, mockAuthorStore)

		output, _ := a.GetBookByID(ctx, v.id)

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc           string
		book           models.Book
		lastInsertedID int
		checkBook      bool
		checkAuthor    bool
		err            error
	}{
		{
			"valid case", models.Book{Title: "RD sharma",
				Author:      models.Author{FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
				Publication: "Arihanth", PublicationDate: "12-08-2011",
			}, 2, false, true, nil,
		},
		{"invalid book already exists", models.Book{
			Title:       "RD",
			Author:      models.Author{FirstName: "gaurav", LastName: "chaudhari", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "Penguin", PublicationDate: "12-08-2011",
		}, 0, true, false, nil,
		},
		{"invalid case error from store layer", models.Book{Title: "xyz",
			Author:      models.Author{FirstName: "gaurav", LastName: "chaudhari", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "NCERT", PublicationDate: "12-08-2011",
		}, 0, false, true, errors.New("err"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		a := New(mockBookStore, mockAuthorStore)
		ctx := context.Background()

		mockBookStore.EXPECT().CheckBook(ctx, &v.book).Return(v.checkBook).AnyTimes()
		mockAuthorStore.EXPECT().CheckAuthorByID(ctx, v.book.Author.ID).Return(v.checkAuthor).AnyTimes()
		mockBookStore.EXPECT().PostBook(ctx, &v.book).Return(v.lastInsertedID, v.err).AnyTimes()

		id, _ := a.PostBook(ctx, &v.book)

		if !reflect.DeepEqual(id, v.lastInsertedID) {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc        string
		book        models.Book
		expOutput   models.Book
		checkAuthor bool
		checkBook   bool
		errAuthor   error
		errBook     error
	}{
		{"valid case", models.Book{ID: 1, Title: "RD sharma",
			Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "12-08-2011",
		}, models.Book{ID: 1, Title: "RD sharma",
			Author:      models.Author{ID: 1, FirstName: "gaurav", LastName: "chandra", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "12-08-2011",
		}, true, true, nil, nil,
		},
		{"invalid book title", models.Book{ID: 2, Title: "",
			Author:      models.Author{ID: 2, FirstName: "gaurav", LastName: "", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "NCERT", PublicationDate: "12-08-2011",
		}, models.Book{}, false, true, nil, nil,
		},
		{"invalid case error from store layer put author", models.Book{ID: 3, Title: "",
			Author:      models.Author{ID: 3, FirstName: "gaurav", LastName: "", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "NCERT", PublicationDate: "12-08-2011",
		}, models.Book{}, true, true,
			errors.New("err"), nil,
		},
		{"invalid case error from store layer put book", models.Book{ID: 4, Title: "",
			Author:      models.Author{ID: 4, FirstName: "gaurav", LastName: "", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "NCERT", PublicationDate: "12-08-2011",
		}, models.Book{}, true, true,
			nil, errors.New("err"),
		},
		{"invalid case book not exist", models.Book{ID: 5, Title: "",
			Author:      models.Author{ID: 5, FirstName: "gaurav", LastName: "", Dob: "18-07-2001", PenName: "GCC"},
			Publication: "NCERT", PublicationDate: "12-08-2011",
		}, models.Book{}, true, false,
			nil, nil,
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		a := New(mockBookStore, mockAuthorStore)
		ctx := context.Background()

		mockAuthorStore.EXPECT().CheckAuthorByID(ctx, v.book.Author.ID).Return(v.checkAuthor).AnyTimes()
		mockAuthorStore.EXPECT().PutAuthor(ctx, v.book.Author.ID, v.book.Author).Return(v.book.Author, v.errAuthor).AnyTimes()
		mockBookStore.EXPECT().CheckBookBid(ctx, v.book.ID).Return(v.checkBook).AnyTimes()
		mockBookStore.EXPECT().PutBook(ctx, v.book.ID, &v.book).Return(v.book, v.errBook).AnyTimes()

		output, _ := a.PutBook(ctx, v.book.ID, &v.book)

		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         int
		rowDeleted int
		checkBook  bool
		err        error
	}{
		{
			"valid case", 1, 1, true, nil,
		},
		{
			"invalid case book not exist", 2, 0, false, nil,
		},
		{
			"invalid case error form store layer deleteBook", 3, 0, true, errors.New("err"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		a := New(mockBookStore, mockAuthorStore)
		ctx := context.Background()

		mockBookStore.EXPECT().CheckBookBid(ctx, v.id).Return(v.checkBook).AnyTimes()
		mockBookStore.EXPECT().DeleteBook(ctx, v.id).Return(v.rowDeleted, v.err).AnyTimes()

		deletedID, _ := a.DeleteBook(ctx, v.id)

		if !reflect.DeepEqual(v.rowDeleted, deletedID) {
			t.Errorf("Expected %v\tGot %v", v.rowDeleted, deletedID)
		}
	}
}

package datastorebook

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"Projects/GoLang-Interns-2022/threeLayer/drivers"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetAllBooks(t *testing.T) {
	testcases := []struct {
		desc           string
		expectedOutput []models.Book
		rows           *sqlmock.Rows
		title          string
		err            error
	}{
		{
			"valid case", []models.Book{
				{ID: 1,
					Title: "RD sharma",
					Author: models.Author{
						ID: 1,
					},
					Publication:     "Arihanth",
					PublicationDate: "12-08-2011",
				},
			}, sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow(1, "RD sharma", "Arihanth", "12-08-2011", 1), "", nil,
		},
		{
			"inValid case", []models.Book{},
			sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow("id", "RD sharma", "Arihanth", "12-08-2011", 1), "", nil,
		},
		{
			"inValid case", []models.Book{},
			sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow(1, "RD sharma", "Arihanth", "12-08-2011", 1), "", errors.New("err"),
		},
	}
	for i, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.SelectFromBook).WillReturnRows(v.rows).WillReturnError(v.err)

		ctx := context.Background()
		output, _ := a.GetAllBooks(ctx)

		if !reflect.DeepEqual(output, v.expectedOutput) {
			fmt.Println(i)
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestGetAllWithTitle(t *testing.T) {
	testcases := []struct {
		desc           string
		expectedOutput []models.Book
		rows           *sqlmock.Rows
		title          string
		err            error
	}{
		{
			"valid case", []models.Book{
				{ID: 1,
					Title: "RD sharma",
					Author: models.Author{
						ID: 1,
					},
					Publication:     "Arihanth",
					PublicationDate: "12-08-2011",
				},
			}, sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow(1, "RD sharma", "Arihanth", "12-08-2011", 1), "RD sharma", nil,
		},
		{
			"inValid case", []models.Book{},
			sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow(1, "RD sharma", "Arihanth", "12-08-2011", 1), "RD sharma", errors.New("err"),
		},
		{
			"inValid case", []models.Book{},
			sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow("id", "RD sharma", "Arihanth", "12-08-2011", 1), "RD sharma", nil,
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.SelectFromBookByTitle).WithArgs(v.title).WillReturnRows(v.rows).WillReturnError(v.err)

		ctx := context.Background()
		output, _ := a.GetAllBooksByTitle(ctx, v.title)

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
		rows           *sqlmock.Rows
		err            error
	}{
		{
			"valid case", 1, models.Book{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, sqlmock.NewRows([]string{"bookID", "title", "publication", "publicationDate",
				"authorID"}).AddRow(1, "RD sharma", "Arihanth", "12-08-2011", 1), nil,
		},

		{"invalid case", 1, models.Book{}, sqlmock.NewRows([]string{"bookID", "title",
			"publication", "publicationDate", "authorID"}), nil},

		{"invalid case", 2, models.Book{}, sqlmock.NewRows([]string{"bookID", "title",
			"publication", "publicationDate", "authorID"}).AddRow("id", "RD sharma", "Arihanth", "12-08-2011", 1), nil},

		{"invalid case", 3, models.Book{}, sqlmock.NewRows([]string{"bookID", "title",
			"publication", "publicationDate", "authorID"}), errors.New("err")},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)
		ctx := context.Background()

		mock.ExpectQuery(drivers.SelectFromBookByID).WithArgs(v.id).WillReturnRows(v.rows).WillReturnError(v.err)

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
		res            driver.Result
		lastInsertedID int
		err            error
	}{
		{
			"valid case", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, sqlmock.NewResult(1, 0), 1, nil,
		},
		{
			"invalid case", models.Book{
				Title: "RD",
				Author: models.Author{
					ID: 42,
				},
				Publication:     "Pengiun",
				PublicationDate: "12-08-2012",
			}, sqlmock.NewErrorResult(errors.New("err")), 0, nil,
		},
		{
			"invalid case", models.Book{
				Title: "RDD",
				Author: models.Author{
					ID: 49,
				},
				Publication:     "Pengiun",
				PublicationDate: "12-08-2012",
			}, sqlmock.NewResult(0, 0), 0, errors.New("err"),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.InsertIntoBook).WithArgs(v.body.Title, v.body.Publication,
			v.body.PublicationDate, v.body.Author.ID).WillReturnResult(v.res).WillReturnError(v.err)

		ctx := context.Background()
		id, _ := a.PostBook(ctx, &v.body)

		if !reflect.DeepEqual(id, v.lastInsertedID) {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         int
		rowDeleted int
		res        driver.Result
		err        error
	}{
		{
			"valid case", 1, 1,
			sqlmock.NewResult(0, 1), nil,
		},
		{
			"invalid case", -1, 0,
			sqlmock.NewResult(0, 0), errors.New("err"),
		},
		{
			"invalid case", -1, 0,
			sqlmock.NewErrorResult(errors.New("err")), nil,
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.DeleteBookQuery).WithArgs(v.id).WillReturnResult(v.res).WillReturnError(v.err)

		ctx := context.Background()

		rowDeleted, _ := a.DeleteBook(ctx, v.id)
		if !reflect.DeepEqual(v.rowDeleted, rowDeleted) {
			t.Errorf("Expected %v\tGot %v", v.rowDeleted, rowDeleted)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		book      models.Book
		expOutput models.Book
		err       error
	}{
		{
			"valid case", 1, models.Book{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, models.Book{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, nil,
		},
		{
			"invalid case", 2, models.Book{
				ID:    2,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, models.Book{}, errors.New("err"),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.UpdateBook).WithArgs(v.book.Title, v.book.Publication, v.book.PublicationDate,
			v.book.ID).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(v.err)

		ctx := context.Background()

		output, _ := a.PutBook(ctx, v.book.ID, &v.book)
		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestDeleteBookByAuthorID(t *testing.T) {
	testcases := []struct {
		desc       string
		id         int
		rowDeleted int
		err        error
	}{
		{
			"valid case", 1, 1, nil,
		},
		{
			"invalid case", -1, 0, errors.New("err"),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.DeleteBookByAuthorID).WithArgs(v.id).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(v.err)

		ctx := context.Background()

		err := a.DeleteBookByAuthorID(ctx, v.id)
		if err != v.err {
			t.Errorf("Expected %v\tGot %v", nil, err)
		}
	}
}

func TestCheckBook(t *testing.T) {
	testcases := []struct {
		desc      string
		book      models.Book
		expOutput bool
		err       error
		row       *sqlmock.Rows
	}{
		{
			"valid case", models.Book{
				ID:    1,
				Title: "RD sharma",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, true, nil, sqlmock.NewRows([]string{"bookID"}).AddRow(1),
		},
		{
			"invalid case", models.Book{
				ID:    2,
				Title: "RD",
				Author: models.Author{
					ID: 1,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, false, errors.New("err"), sqlmock.NewRows([]string{"bookID"}),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.CheckBook).WithArgs(v.book.Title, v.book.Publication, v.book.PublicationDate,
			v.book.ID).WillReturnRows(v.row).WillReturnError(v.err)

		ctx := context.Background()

		output := a.CheckBook(ctx, &v.book)
		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestCheckBookBid(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		expOutput bool
		err       error
		row       *sqlmock.Rows
	}{
		{
			"valid case", 1,
			true, nil, sqlmock.NewRows([]string{"bookID"}).AddRow(1),
		},
		{
			"invalid case", 2,
			false, errors.New("err"), sqlmock.NewRows([]string{"bookID"}),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.CheckBookBYID).WithArgs(v.id).WillReturnRows(v.row).WillReturnError(v.err)

		ctx := context.Background()

		output := a.CheckBookBid(ctx, v.id)
		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

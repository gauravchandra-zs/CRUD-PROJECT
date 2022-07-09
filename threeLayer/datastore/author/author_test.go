package datastoreauthor

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
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

func TestAuthorStore_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		id        int
		author    models.Author
		expOutput models.Author
		rows      *sqlmock.Rows
		err       error
	}{
		{
			"author present", 1, models.Author{
				ID:        1,
				FirstName: "GAURAV",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{
				ID:        1,
				FirstName: "GAURAV",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, sqlmock.NewRows([]string{"authorID", "firstName", "lastName", "dob",
				"penName"}).AddRow(41, "Gaurav", "chandra", "18-07-2001", "GCC"), nil,
		},
		{
			"author not present", 23, models.Author{
				ID:        23,
				FirstName: "xyz",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
			sqlmock.NewRows([]string{"authorID", "firstName", "lastName", "dob", "penName"}), errors.New("err"),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.UpdateAuthor).WithArgs(v.author.FirstName, v.author.LastName,
			v.author.Dob, v.author.PenName, v.id).WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(v.err)

		ctx := context.Background()
		output, _ := a.PutAuthor(ctx, v.id, v.author)

		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestAuthorStore_PostAuthor(t *testing.T) {
	tc := []struct {
		desc           string
		body           models.Author
		lastInsertedID int
		result         driver.Result
		err            error
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 1, sqlmock.NewResult(1, 1), nil,
		},
		{
			"invalid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0, sqlmock.NewResult(0, 0), errors.New("error"),
		},
		{
			"invalid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0, sqlmock.NewErrorResult(errors.New("error")), nil,
		},
	}
	for _, v := range tc {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.InsertIntoAuthor).WithArgs(v.body.FirstName, v.body.LastName,
			v.body.Dob, v.body.PenName).WillReturnResult(v.result).WillReturnError(v.err)

		ctx := context.Background()
		id, _ := a.PostAuthor(ctx, v.body)
		if !reflect.DeepEqual(id, v.lastInsertedID) {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}

}

func TestAuthorStore_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		id         int
		rowDeleted int
		err        error
		res        driver.Result
	}{
		{
			"valid case", 1,
			1, nil, sqlmock.NewResult(0, 1),
		},
		{
			"invalid case", -1,
			0, nil, sqlmock.NewErrorResult(errors.New("err")),
		},
		{
			"invalid case", -1,
			0, errors.New("err"), sqlmock.NewResult(0, 0),
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectExec(drivers.DeleteAuthorQuery).WithArgs(v.id).WillReturnResult(v.res).WillReturnError(v.err)

		ctx := context.Background()
		id, _ := a.DeleteAuthor(ctx, v.id)

		if !reflect.DeepEqual(id, v.rowDeleted) {
			t.Errorf("Expected %v\tGot %v", v.rowDeleted, id)
		}
	}
}

func TestAuthorStore_GetAuthorByID(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput models.Author
		rows           *sqlmock.Rows
		err            error
	}{
		{
			"author present", 1, models.Author{
				ID:        1,
				FirstName: "RD",
				LastName:  "sharma",
				Dob:       "07-10-1998",
				PenName:   "GCC",
			}, sqlmock.NewRows([]string{"authorID", "firstName", "lastName", "dob",
				"penName"}).AddRow(1, "RD", "sharma", "07-10-1998", "GCC"), nil,
		},
		{
			"invalid case", 23,
			models.Author{}, sqlmock.NewRows([]string{"authorID", "firstName",
				"lastName", "dob", "penName"}), nil,
		},
		{
			"invalid case", 50,
			models.Author{}, sqlmock.NewRows([]string{"authorID", "firstName",
				"lastName", "dob", "penName"}), errors.New("err"),
		},
		{
			"invalid case", 1, models.Author{},
			sqlmock.NewRows([]string{"authorID", "firstName", "lastName", "dob",
				"penName"}).AddRow("id", "RD", "sharma", "07-10-1998", "GCC"), nil,
		},
	}
	for _, v := range testcases {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.SelectAuthorByID).WithArgs(v.id).WillReturnRows(v.rows).WillReturnError(v.err)
		ctx := context.Background()
		output, _ := a.GetAuthorByID(ctx, v.id)
		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("expected %v,got %v", v.expectedOutput, output)
		}
	}
}

func TestAuthorStore_CheckAuthor(t *testing.T) {
	tc := []struct {
		desc      string
		body      models.Author
		row       *sqlmock.Rows
		err       error
		expOutput bool
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, sqlmock.NewRows([]string{"authorID"}).AddRow(1), nil, true,
		},
		{
			"invalid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, sqlmock.NewRows([]string{"authorID"}), errors.New("error"), false,
		},
	}
	for _, v := range tc {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.CheckAuthor).WithArgs(v.body.FirstName, v.body.LastName,
			v.body.Dob, v.body.PenName).WillReturnRows(v.row).WillReturnError(v.err)

		ctx := context.Background()
		flag := a.CheckAuthor(ctx, v.body)
		if !reflect.DeepEqual(flag, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, flag)
		}
	}

}

func TestAuthorStore_CheckAuthorByID(t *testing.T) {
	tc := []struct {
		desc      string
		id        int
		row       *sqlmock.Rows
		err       error
		expOutput bool
	}{
		{
			"valid case", 1,
			sqlmock.NewRows([]string{"authorID"}).AddRow(1), nil, true,
		},
		{
			"invalid case", 2,
			sqlmock.NewRows([]string{"authorID"}), errors.New("error"), false,
		},
	}
	for _, v := range tc {
		db, mock := NewMock()
		a := New(db)

		mock.ExpectQuery(drivers.CheckAuthorBYID).WithArgs(v.id).WillReturnRows(v.row).WillReturnError(v.err)

		ctx := context.Background()
		flag := a.CheckAuthorByID(ctx, v.id)
		if !reflect.DeepEqual(flag, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, flag)
		}
	}

}

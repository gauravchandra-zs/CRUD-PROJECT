package datastorebook

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"Projects/GoLang-Interns-2022/threeLayer/models"

	_ "github.com/go-sql-driver/mysql"
)

func initializeMySQL(t *testing.T) *sql.DB {
	var Driver = "mysql"

	var DataSource = "root:@Gc18072001@/test"

	db, err := sql.Open(Driver, DataSource)
	if err != nil {
		t.Errorf("could not connect to sql, err:%v", err)
	}

	return db
}
func TestBookStore_PostBook(t *testing.T) {
	testcases := []struct {
		desc           string
		body           models.Book
		lastInsertedID int
	}{
		{
			"Book exist", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, 0,
		},
		{
			"Book not exist", models.Book{
				Title: "RD",
				Author: models.Author{
					ID: 42,
				},
				Publication:     "Pengiun",
				PublicationDate: "12-08-2012",
			}, 0,
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		id, err := a.PostBook(&v.body)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(id, v.lastInsertedID) {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}
}

func TestBookStore_GetAll(t *testing.T) {
	testcases := []struct {
		desc           string
		expectedOutput []models.Book
		title          string
	}{
		{
			"valid case", []models.Book{
				{ID: 59,
					Title: "RD sharma",
					Author: models.Author{
						ID: 40,
					},
					Publication:     "Arihanth",
					PublicationDate: "12-08-2011",
				},
			}, "",
		},
		{
			"valid case", []models.Book{
				{ID: 59,
					Title: "RD sharma",
					Author: models.Author{
						ID: 40,
					},
					Publication:     "Arihanth",
					PublicationDate: "12-08-2011",
				},
			}, "RD sharma",
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		output, err := a.GetAll(v.title)
		if err != nil {
			t.Errorf("test case fail getAll returns error")
		}

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}
func TestBookStore_GetByID(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		expectedOutput models.Book
	}{
		{
			"valid case", 59, models.Book{
				ID:    59,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		},
		{"invalid case", 1, models.Book{}},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		output, err := a.GetByID(v.id)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestBookStore_PutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		body      models.Book
		expOutput models.Book
	}{
		{
			"valid case", models.Book{
				ID:    59,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, models.Book{
				ID:    59,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		},
		{
			"invalid case", models.Book{
				ID:    60,
				Title: "RD sharma",
				Author: models.Author{
					ID: 40,
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, models.Book{},
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		output, err := a.PutBook(&v.body)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(output, v.expOutput) {
			t.Errorf("Expected %v\tGot %v", v.expOutput, output)
		}
	}
}

func TestBookStore_DeleteBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         int
		rowDeleted int
	}{
		{
			"valid case", 1, 0,
		},
		{
			"invalid case", -1, 0,
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		deletedID, err := a.DeleteBook(v.id)
		if err != nil {
			t.Errorf("test case fail")
		}

		if !reflect.DeepEqual(v.rowDeleted, deletedID) {
			t.Errorf("Expected %v\tGot %v", v.rowDeleted, deletedID)
		}
	}
}

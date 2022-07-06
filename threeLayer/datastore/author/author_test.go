package datastoreauthor

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

func TestAuthorStore_PostAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		body           models.Author
		lastInsertedID int
	}{
		{
			"author exist", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 1,
		},
		{
			"author not exist", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 1,
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		id, err := a.PostAuthor(v.body)
		if err != nil {
			t.Errorf("test case fail ,error return by PostBook")
		}

		if id == 0 {
			t.Errorf("Expected %v\tGot %v", v.lastInsertedID, id)
		}
	}
}

func TestAuthorStore_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		author         models.Author
		expectedOutput models.Author
	}{
		{
			"author present", 41, models.Author{
				ID:        41,
				FirstName: "GAURAV",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{
				ID:        41,
				FirstName: "GAURAV",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			},
		},
		{
			"author not present", 23, models.Author{
				ID:        23,
				FirstName: "xyz",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		output, err := a.PutAuthor(v.id, v.author)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("expected %v,got %v", v.expectedOutput, output)
		}
	}
}

func TestAuthorStore_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc          string
		id            int
		expRowDeleted int
	}{
		{
			"valid case", 1, 0,
		},
		{
			"valid case", -1, 0,
		},
	}
	for _, v := range testcases {
		db := initializeMySQL(t)
		a := New(db)

		rowDeleted, err := a.DeleteAuthor(v.id)
		if err != nil {
			log.Print(err)
		}

		if rowDeleted != v.expRowDeleted {
			t.Errorf("test case fail ,expected %v got %v", v.expRowDeleted, rowDeleted)
		}
	}
}

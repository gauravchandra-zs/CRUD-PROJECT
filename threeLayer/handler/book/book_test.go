package handlerbook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"Projects/GoLang-Interns-2022/threeLayer/models"

	"github.com/gorilla/mux"
)

type mockService struct{}

func (m mockService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return []models.Book{
		{
			ID:              1,
			Title:           "RD sharma",
			Author:          models.Author{},
			Publication:     "Arihanth",
			PublicationDate: "12-08-2011",
		}}, nil
}

func (m mockService) GetBookByID(ctx context.Context, id int) (models.Book, error) {
	if id <= 0 {
		return models.Book{}, errors.New("invalid id")
	}

	return models.Book{
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
	}, nil
}

func (m mockService) PostBook(ctx context.Context, book *models.Book) (int, error) {
	if !validateBook(book) || !validateAuthor(book.Author) {
		return 0, errors.New("invalid book or author")
	}

	return 1, nil
}

func (m mockService) DeleteBook(ctx context.Context, id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	return 1, nil
}

func (m mockService) PutBook(ctx context.Context, id int, book *models.Book) (models.Book, error) {
	if book.ID <= 0 || !validateBook(book) || !validateAuthor(book.Author) {
		return models.Book{}, errors.New("invalid book or author")
	}

	return models.Book{
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
	}, nil
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc   string
		body   models.Book
		status int
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
			}, http.StatusCreated,
		},
		{
			"invalid case", models.Book{
				Title: "",
				Author: models.Author{
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "NCERT",
				PublicationDate: "12-08-2011",
			}, http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(myData))
		w := httptest.NewRecorder()
		a := New(mockService{})

		a.PostBook(w, req)

		if !reflect.DeepEqual(w.Result().StatusCode, v.status) {
			t.Errorf("Expected %v\tGot %v", v.status, w.Result().StatusCode)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		expectedStatus int
	}{
		{
			"valid case", "1", http.StatusNoContent,
		},
		{
			"valid case", "-1", http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/books/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		w := httptest.NewRecorder()
		a := New(mockService{})

		a.DeleteBook(w, req)

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) {
			t.Errorf("Expected %v\tGot %v", v.expectedStatus, w.Result().StatusCode)
		}
	}
}

func TestGetAllBooks(t *testing.T) {
	testcases := []struct {
		desc           string
		expectedOutput []models.Book
	}{
		{"valid case", []models.Book{
			{
				ID:              1,
				Title:           "RD sharma",
				Author:          models.Author{},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		},
		},
	}
	for _, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()

		a := New(mockService{})

		a.GetAllBooks(w, req)

		data, err := io.ReadAll(w.Body)
		if err != nil {
			t.Errorf("test case fail ,error in reading body")
		}

		var output []models.Book

		err = json.Unmarshal(data, &output)
		if err != nil {
			t.Errorf("test case fail ,error in unmarshaling data")
		}

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestGetBookByID(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		expectedOutput models.Book
		expectedStatus int
	}{
		{"valid case", "1", models.Book{
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
		}, http.StatusOK,
		},
		{"invalid case", "-1", models.Book{}, http.StatusBadRequest},
	}
	for _, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/book/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		w := httptest.NewRecorder()
		a := New(mockService{})

		a.GetBookById(w, req)

		data, err := io.ReadAll(w.Body)
		if err != nil {
			log.Print(err)
		}

		var output models.Book

		err = json.Unmarshal(data, &output)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(output, v.expectedOutput) && !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc           string
		body           models.Book
		expectedOutput models.Book
		expectedStatus int
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
			}, http.StatusAccepted,
		},

		{
			"invalid case", models.Book{
				ID:    2,
				Title: "",
				Author: models.Author{
					ID:        1,
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "NCERT",
				PublicationDate: "12-08-2011",
			}, models.Book{},
			http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf(" error in marshaling")
		}

		req := httptest.NewRequest(http.MethodPut, "/books/{id}", bytes.NewBuffer(myData))
		a := New(mockService{})
		w := httptest.NewRecorder()

		a.PutBook(w, req)

		myData, err = io.ReadAll(w.Body)
		if err != nil {
			log.Print(err)
		}

		var output models.Book

		err = json.Unmarshal(myData, &output)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) && !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func validateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}
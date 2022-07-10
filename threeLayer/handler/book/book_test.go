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
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"
)

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc   string
		body   any
		status int
		err    error
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
			}, http.StatusCreated, nil,
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
			}, http.StatusBadRequest, nil,
		},
		{
			"invalid case", "Book", http.StatusInternalServerError, nil,
		},
		{
			"invalid case", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			}, http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "XYZ",
				PublicationDate: "12-08-2011",
			}, http.StatusBadRequest, nil,
		},
		{
			"invalid case", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Arihanth",
				PublicationDate: "2000",
			}, http.StatusBadRequest, nil,
		},
		{
			"invalid case", models.Book{
				Title: "RD sharma",
				Author: models.Author{
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Arihanth",
				PublicationDate: "18-07-2023",
			}, http.StatusBadRequest, nil,
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
				Publication:     "Arihanth",
				PublicationDate: "18-07-2000",
			}, http.StatusBadRequest, nil,
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockServiceBook := service.NewMockBook(mockCtrl)
		a := New(mockServiceBook)

		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(myData))
		w := httptest.NewRecorder()
		ctx := context.Background()
		book, _ := v.body.(models.Book)

		mockServiceBook.EXPECT().PostBook(ctx, &book).Return(1, v.err).AnyTimes()

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
		err            error
	}{
		{
			"valid case", "1", http.StatusNoContent, nil,
		},
		{
			"valid case", "2", http.StatusBadRequest, errors.New("err"),
		},
		{
			"valid case", "-1", http.StatusBadRequest, nil,
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockServiceBook := service.NewMockBook(mockCtrl)
		a := New(mockServiceBook)

		req := httptest.NewRequest(http.MethodGet, "/book/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		id, _ := strconv.Atoi(v.id)
		w := httptest.NewRecorder()

		mockServiceBook.EXPECT().DeleteBook(context.Background(), id).Return(1, v.err).AnyTimes()
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
		err            error
	}{
		{"valid case", []models.Book{
			{
				ID:              1,
				Title:           "RD sharma",
				Author:          models.Author{},
				Publication:     "Arihanth",
				PublicationDate: "12-08-2011",
			},
		}, nil,
		},
		{"invalid case", []models.Book{}, errors.New("err")},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockServiceBook := service.NewMockBook(mockCtrl)
		a := New(mockServiceBook)
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()

		ctx := context.Background()
		ctx = context.WithValue(ctx, "title", "")
		ctx = context.WithValue(ctx, "includeAuthor", "")

		mockServiceBook.EXPECT().GetAllBooks(ctx).Return(v.expectedOutput, v.err).AnyTimes()

		a.GetAllBooks(w, req)

		data, err := io.ReadAll(w.Body)
		if err != nil {
			t.Errorf("test case fail ,error in reading body")
		}

		var output []models.Book

		err = json.Unmarshal(data, &output)
		if err != nil {
			log.Print("test case fail ,error in unmarshaling data")
			continue
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
		err            error
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
		}, http.StatusOK, nil,
		},
		{"invalid case", "-1", models.Book{},
			http.StatusBadRequest, nil},
		{"invalid case", "2", models.Book{},
			http.StatusBadRequest, errors.New("err")},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockServiceBook := service.NewMockBook(mockCtrl)
		a := New(mockServiceBook)
		req := httptest.NewRequest(http.MethodGet, "/book/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		w := httptest.NewRecorder()
		id, _ := strconv.Atoi(v.id)

		mockServiceBook.EXPECT().GetBookByID(context.Background(), id).Return(v.expectedOutput, v.err).AnyTimes()

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
		id             int
		body           any
		expectedOutput models.Book
		expectedStatus int
		err            error
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
			}, http.StatusAccepted, nil,
		},

		{
			"invalid case", 2, models.Book{
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
			http.StatusBadRequest, nil,
		},
		{
			"invalid case", 3, models.Book{
				ID:    3,
				Title: "xyz",
				Author: models.Author{
					ID:        1,
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Scholastic",
				PublicationDate: "12-08-2011",
			}, models.Book{},
			http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case", -3, models.Book{
				ID:    -3,
				Title: "xyz",
				Author: models.Author{
					ID:        1,
					FirstName: "gaurav",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Scholastic",
				PublicationDate: "12-08-2011",
			}, models.Book{},
			http.StatusBadRequest, nil,
		},
		{
			"invalid case", 3, models.Book{
				ID:    3,
				Title: "xyz",
				Author: models.Author{
					ID:        1,
					FirstName: "",
					LastName:  "chandra",
					Dob:       "18-07-2001",
					PenName:   "GCC",
				},
				Publication:     "Scholastic",
				PublicationDate: "12-08-2011",
			}, models.Book{},
			http.StatusBadRequest, nil,
		},
		{
			"invalid case", 4, "book", models.Book{},
			http.StatusBadRequest, nil,
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockServiceBook := service.NewMockBook(mockCtrl)
		a := New(mockServiceBook)
		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf(" error in marshaling")
		}

		req := httptest.NewRequest(http.MethodPut, "/book/{id}", bytes.NewBuffer(myData))
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(v.id)})
		book, _ := v.body.(models.Book)
		mockServiceBook.EXPECT().PutBook(context.Background(), v.id, &book).Return(v.expectedOutput, v.err).AnyTimes()

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

package handlerauthor

import (
	"bytes"
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

type mockService struct {
}

func (m mockService) PostAuthor(author models.Author) (int, error) {
	if !ValidateAuthor(author) {
		return 0, errors.New("invalid author detail")
	}

	return 1, nil
}

func (m mockService) DeleteAuthor(id int) (int, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}
	if id == 100 {
		return 0, errors.New("author not exist")
	}
	return 1, nil
}

func (m mockService) PutAuthor(id int, author models.Author) (models.Author, error) {
	if !ValidateAuthor(author) || id <= 0 {
		return models.Author{}, errors.New("invalid author detail")
	}

	author.ID = 1

	return author, nil
}
func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		body           models.Author
		expectedStatus int
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusCreated,
		},
		{
			"invalid case", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPost, "/author", bytes.NewBuffer(myData))
		w := httptest.NewRecorder()
		a := New(mockService{})

		a.PostAuthor(w, req)

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) {
			t.Errorf("Expected %v\tGot %v", v.expectedStatus, w.Result().StatusCode)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		body           models.Author
		expectedOutput models.Author
		expectedStatus int
	}{
		{"valid case", "1", models.Author{
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
		}, http.StatusCreated},
		{
			"invalid case", "1", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, http.StatusBadRequest,
		},
		{
			"invalid case", "-1", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPut, "/author", bytes.NewBuffer(myData))
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		a := New(mockService{})
		w := httptest.NewRecorder()

		a.PutAuthor(w, req)

		data, err := io.ReadAll(w.Body)
		if err != nil {
			log.Print(err)
		}

		var output models.Author

		err = json.Unmarshal(data, &output)
		if err != nil {
			log.Print(err)
		}

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) && !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		id             string
		expectedStatus int
	}{
		{
			"valid case", "1", http.StatusNoContent,
		},
		{
			"invalid case", "100", http.StatusBadRequest,
		},
		{
			"invalid case", "-1", http.StatusBadRequest,
		},
	}
	for _, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/author/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		w := httptest.NewRecorder()
		a := New(mockService{})

		a.DeleteAuthor(w, req)

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) {
			t.Errorf("Expected %v\tGot %v", v.expectedStatus, w.Result().StatusCode)
		}
	}
}

func ValidateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

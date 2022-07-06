package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc               string
		input              string
		expectedOutput     []Book
		expectedStatusCode int
	}{
		{desc: "valid url.", input: "http://localhost:8000/books", expectedOutput: []Book{
			{ID: 32, Title: "title1", Author: Author{},
				Publication: "Arihanth", PublicationDate: "18-08-2018"},
		}, expectedStatusCode: http.StatusOK,
		},
		{desc: "valid url with book title", input: "http://localhost:8000/books?title=title1", expectedOutput: []Book{
			{ID: 32, Title: "title1", Author: Author{},
				Publication: "Arihanth", PublicationDate: "18-08-2018"},
		}, expectedStatusCode: http.StatusOK,
		},
		{desc: "valid url with includeAuthor = true ", input: "http://localhost:8000/books?includeAuthor=true", expectedOutput: []Book{
			{ID: 32, Title: "title1", Author: Author{ID: 4, FirstName: "gaurav1", LastName: "chandra1", Dob: "07-10-1998", PenName: "GCC"},
				Publication: "Arihanth", PublicationDate: "18-08-2018"},
		}, expectedStatusCode: http.StatusOK,
		},
		{desc: "valid url with includeAuthor = false", input: "http://localhost:8000/books?includeAuthor=false", expectedOutput: []Book{
			{ID: 32, Title: "title1", Publication: "Arihanth", PublicationDate: "18-08-2018"},
			{ID: 33, Title: "title1", Publication: "Arihanth", PublicationDate: "18-08-2018"},
		}, expectedStatusCode: http.StatusOK,
		},
	}
	for i, v := range testcases {
		newRequest := httptest.NewRequest(http.MethodGet, v.input, nil)
		response := httptest.NewRecorder()
		GetAll(response, newRequest)

		if response.Result().StatusCode != v.expectedStatusCode {
			t.Errorf("status code not matched test case fail at %d", i)
		}

		var output []Book

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Print("error in reading data from response.Body")
		}

		err = json.Unmarshal(body, &output)
		if err != nil {
			return
		}

		fmt.Println(output)

		if !reflect.DeepEqual(output, v.expectedOutput) {
			t.Errorf("test case fail at %d", i)
		}
	}
}

func TestGetByID(t *testing.T) {
	testcases := []struct {
		desc               string
		input              string
		expectedOutput     Book
		expectedStatusCode int
	}{
		{desc: "test for valid case: ", input: "35", expectedOutput: Book{ID: 32, Title: "title1",
			Author:      Author{ID: 4, FirstName: "gaurav1", LastName: "chandra1", Dob: "07-10-1998", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "18-08-2018"},
			expectedStatusCode: http.StatusOK},
		{"test for invalid id : ", "-1",
			Book{}, http.StatusBadRequest},
		{"test for valid id but id not exist : ", "9440",
			Book{}, http.StatusNotFound},
	}
	for i, v := range testcases {
		newRequest := httptest.NewRequest(http.MethodGet, "http://localhost:8000/books/{id}"+v.input, nil)
		response := httptest.NewRecorder()
		newRequest = mux.SetURLVars(newRequest, map[string]string{"id": v.input})
		GetByID(response, newRequest)

		if response.Result().StatusCode != v.expectedStatusCode {
			t.Errorf("test cases fail")
		} else {
			var output Book
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Print("error in reading data from response.Body")
			}
			err = json.Unmarshal(body, &output)
			if err != nil {
				log.Print(err)
				return
			}
			if !reflect.DeepEqual(output, v.expectedOutput) {
				t.Errorf("test case fail at %d", i)
			}
		}
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc       string
		inputData  Book
		statusCode int
	}{
		{desc: "Details posted.", inputData: Book{Title: "title1", Author: Author{FirstName: "gaurav1",
			LastName: "chandra1", Dob: "07-10-1998", PenName: "GCC"}, Publication: "Arihanth", PublicationDate: "18-08-2018"},
			statusCode: http.StatusCreated},
		{desc: "Details posted.", inputData: Book{Title: "title1", Author: Author{FirstName: "gaurav",
			LastName: "chandra1", Dob: "07-10-1998", PenName: "GCC"}, Publication: "Arihanth",
			PublicationDate: "18-08-2018"},
			statusCode: http.StatusCreated},
		{desc: "Invalid book name.", inputData: Book{Title: "", Author: Author{FirstName: "gaurav2",
			LastName: "chandra2", Dob: "07-10-1998", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "21-04-1985"},
			statusCode: http.StatusBadRequest},
		{desc: "Invalid author.", inputData: Book{Title: "title3", Author: Author{FirstName: "",
			LastName: "chandra3", Dob: "19-12-1972", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "21-04-1985"},
			statusCode: http.StatusBadRequest,
		},
		{desc: "Invalid publication date.", inputData: Book{Title: "title4", Author: Author{FirstName: "gaurav4",
			LastName: "chandra4", Dob: "19-12-1972", PenName: "GCC"},
			Publication: "Pengiun", PublicationDate: "21-04-1879"},
			statusCode: http.StatusBadRequest,
		},
		{desc: "Invalid publication date.", inputData: Book{Title: "title5",
			Author:      Author{FirstName: "gaurav5", LastName: "chandra5", Dob: "19-12-1972", PenName: "GCC"},
			Publication: "Arihanth", PublicationDate: "21-04-2031"},
			statusCode: http.StatusBadRequest,
		},
	}

	for i, v := range testcases {
		myData, err := json.Marshal(v.inputData)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		newRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8000/book", bytes.NewBuffer(myData))
		response := httptest.NewRecorder()
		PostBook(response, newRequest)

		if !reflect.DeepEqual(response.Result().StatusCode, v.statusCode) {
			fmt.Println(response.Result().StatusCode)
			t.Errorf("test cases fail at %d", i)
		}
	}
}
func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		inputData  Author
		statusCode int
	}{
		{desc: "Valid case.", inputData: Author{FirstName: "Gaurav", LastName: "chandra", Dob: "01-07-2001", PenName: "GCC"},
			statusCode: http.StatusCreated},
		{desc: "Invalid case.", inputData: Author{FirstName: "", LastName: "chandra", Dob: "81-07-2001", PenName: "GCC"},
			statusCode: http.StatusBadRequest},
		{desc: "Invalid case.", inputData: Author{FirstName: "gaurav1", LastName: "", Dob: "81-07-2001", PenName: "GCC"},
			statusCode: http.StatusBadRequest},
	}

	for i, v := range testcases {
		myData, err := json.Marshal(v.inputData)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		newRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8000/author", bytes.NewReader(myData))
		response := httptest.NewRecorder()

		PostAuthor(response, newRequest)

		if !reflect.DeepEqual(response.Result().StatusCode, v.statusCode) {
			t.Errorf("test cases fail at %d", i)
		}
	}
}
func TestDelete(t *testing.T) {
	testcases := []struct {
		desc       string
		ID         string
		url        string
		statusCode int
	}{
		{desc: "Valid case author.", ID: "26", url: "http://localhost:8000/author/{id}",
			statusCode: http.StatusBadRequest},
		{desc: "Invalid case author", ID: "-1", url: "http://localhost:8000/author/{id}",
			statusCode: http.StatusBadRequest},
		{desc: "valid case book", ID: "32", url: "http://localhost:8000/books/{id}",
			statusCode: http.StatusBadRequest},
		{desc: "Invalid case book.", ID: "-2", url: "http://localhost:8000/books/{id}",
			statusCode: http.StatusBadRequest},
	}

	for i, v := range testcases {
		newRequest := httptest.NewRequest(http.MethodDelete, v.url+v.ID, nil)
		response := httptest.NewRecorder()
		newRequest = mux.SetURLVars(newRequest, map[string]string{"id": v.ID})
		DeleteAuthor(response, newRequest)

		if !reflect.DeepEqual(response.Result().StatusCode, v.statusCode) {
			fmt.Println(response.Result().StatusCode)
			t.Errorf("test cases fail at %d", i)
		}
	}
}
func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		id         string
		putData    Author
		statusCode int
	}{
		{desc: "Valid case update firstname.", id: "20", putData: Author{ID: 20, FirstName: "mehjndj",
			LastName: "chandra", Dob: "01-07-2001", PenName: "GCC"},
			statusCode: http.StatusCreated},

		{desc: "Valid case id not present.", id: "1000", putData: Author{ID: 1000,
			FirstName: "mehul", LastName: "chandra", Dob: "01-07-2001", PenName: "GCC"},
			statusCode: http.StatusBadRequest},
	}
	for i, v := range testcases {
		myData, err := json.Marshal(v.putData)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		newRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8000/author/{id}"+v.id, bytes.NewReader(myData))
		response := httptest.NewRecorder()
		newRequest = mux.SetURLVars(newRequest, map[string]string{"id": v.id})

		PutAuthor(response, newRequest)

		if !reflect.DeepEqual(response.Result().StatusCode, v.statusCode) {
			t.Errorf("test cases fail at %d", i)
		}
	}
}
func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc       string
		id         string
		inputData  Book
		statusCode int
	}{
		{desc: "valid case id exist", id: "35", inputData: Book{ID: 35, Title: "titl",
			Author: Author{ID: 16, FirstName: "grav",
				LastName: "chan1", Dob: "07-10-1998", PenName: "GCC"}, Publication: "Arihanth",
			PublicationDate: "18-08-2018"},
			statusCode: http.StatusOK},
		{desc: "invalid case id not exist", id: "1000", inputData: Book{ID: 1000, Title: "title1",
			Author: Author{ID: 9, FirstName: "gaurav",
				LastName: "chandra1", Dob: "07-10-1998", PenName: "GCC"}, Publication: "Arihanth",
			PublicationDate: "18-08-2018"},
			statusCode: http.StatusBadRequest},
		{desc: "Invalid book name.", id: "2", inputData: Book{ID: 2, Title: "",
			Author:      Author{ID: 2, FirstName: "gaurav2", LastName: "chandra2", Dob: "07-10-1998", PenName: "GCC"},
			Publication: "Oxford", PublicationDate: "21-04-1985"},
			statusCode: http.StatusBadRequest},
		{desc: "Invalid author.", id: "8", inputData: Book{ID: 8, Title: "title3",
			Author:      Author{ID: 2, FirstName: "", LastName: "chandra3", Dob: "19-12-1972", PenName: "GCC"},
			Publication: "Oxford", PublicationDate: "21-04-1985"},
			statusCode: http.StatusBadRequest,
		},
	}
	for i, v := range testcases {
		myData, err := json.Marshal(v.inputData)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		newRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8000/books/{id}"+v.id, bytes.NewBuffer(myData))
		response := httptest.NewRecorder()
		newRequest = mux.SetURLVars(newRequest, map[string]string{"id": v.id})

		PutBook(response, newRequest)

		if !reflect.DeepEqual(response.Result().StatusCode, v.statusCode) {
			fmt.Println(response.Result().StatusCode)
			t.Errorf("test cases fail at %d", i)
		}
	}
}

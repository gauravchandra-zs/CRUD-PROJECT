package handlerauthor

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

	"Projects/GoLang-Interns-2022/threeLayer/datastore"
	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"
)

func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		body           any
		expectedStatus int
		err            error
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusCreated, nil,
		},
		{
			"invalid author detail", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusBadRequest, nil,
		},
		{
			"invalid body", "hhds", http.StatusBadRequest, nil,
		},
		{
			"error from service layer", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusBadRequest, errors.New("err"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockServiceAuthor := service.NewMockAuthor(mockCtrl)
		a := New(mockServiceAuthor)

		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPost, "/author", bytes.NewBuffer(myData))
		w := httptest.NewRecorder()
		ctx := context.Background()
		mockServiceAuthor.EXPECT().PostAuthor(ctx, v.body).Return(1, v.err).AnyTimes()

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
		body           any
		expectedOutput models.Author
		expectedStatus int
		err            error
	}{
		{
			"valid case", "1", models.Author{
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
			}, http.StatusCreated, nil,
		},
		{
			"invalid author detail", "1", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, http.StatusBadRequest, nil,
		},
		{
			"invalid body type not Author", "2", "author",
			models.Author{}, http.StatusBadRequest, nil,
		},
		{
			"invalid case error from service layer", "1", models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
			http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case id negative ", "-2", "author",
			models.Author{}, http.StatusBadRequest, nil,
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockServicAuthor := service.NewMockAuthor(mockCtrl)
		a := New(mockServicAuthor)

		myData, err := json.Marshal(v.body)
		if err != nil {
			t.Errorf("can not convert data into []byte")
		}

		req := httptest.NewRequest(http.MethodPut, "/author", bytes.NewBuffer(myData))
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		w := httptest.NewRecorder()

		ctx := context.Background()
		id, _ := strconv.Atoi(v.id)
		mockServicAuthor.EXPECT().PutAuthor(ctx, id, v.body).Return(v.expectedOutput, v.err).AnyTimes()

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
		err            error
	}{
		{
			"valid case", "1", http.StatusNoContent, nil,
		},
		{
			"invalid case error from service layer", "100", http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case negative id", "-1", http.StatusBadRequest, nil,
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockServiceAuthor := datastore.NewMockAuthor(mockCtrl)
		a := New(mockServiceAuthor)

		req := httptest.NewRequest(http.MethodGet, "/author/{id}"+v.id, nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		id, _ := strconv.Atoi(v.id)

		w := httptest.NewRecorder()
		ctx := context.Background()
		mockServiceAuthor.EXPECT().DeleteAuthor(ctx, id).Return(1, v.err).AnyTimes()
		a.DeleteAuthor(w, req)

		if !reflect.DeepEqual(w.Result().StatusCode, v.expectedStatus) {
			t.Errorf("Expected %v\tGot %v", v.expectedStatus, w.Result().StatusCode)
		}
	}
}

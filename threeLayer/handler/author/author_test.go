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
			"invalid case", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusBadRequest, nil,
		},
		{
			"invalid case", "hhds", http.StatusBadRequest, nil,
		},
		{
			"invalid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, http.StatusBadRequest, errors.New("err"),
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

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
			"invalid case", "1", models.Author{
				FirstName: "",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, http.StatusBadRequest, nil,
		},
		{
			"invalid case", "2", "author",
			models.Author{}, http.StatusBadRequest, nil,
		},
		{
			"valid case", "1", models.Author{
				ID:        1,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{},
			http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case", "-2", "author",
			models.Author{}, http.StatusBadRequest, nil,
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

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
			"invalid case", "100", http.StatusBadRequest, errors.New("err"),
		},
		{
			"invalid case", "-1", http.StatusBadRequest, nil,
		},
	}
	for _, v := range testcases {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

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

package serviceauthor

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"Projects/GoLang-Interns-2022/threeLayer/datastore"
	"Projects/GoLang-Interns-2022/threeLayer/models"
)

func TestServiceAuthor_PostAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		author      models.Author
		id          int
		checkAuthor bool
		err         error
	}{
		{
			"valid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 1, false, nil,
		},
		{
			"invalid case", models.Author{
				FirstName: "xyz",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0, true, nil,
		},
		{
			"invalid case", models.Author{
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, 0, false, errors.New("error"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		ctx := context.Background()
		a := New(mockAuthorStore, mockBookStore)

		mockAuthorStore.EXPECT().CheckAuthor(ctx, v.author).Return(v.checkAuthor).AnyTimes()
		mockAuthorStore.EXPECT().PostAuthor(ctx, v.author).Return(v.id, v.err).AnyTimes()
		id, _ := a.PostAuthor(ctx, v.author)

		if !reflect.DeepEqual(id, v.id) {
			t.Errorf("Expected %v\tGot %v", v.id, id)
		}
	}
}

func TestServiceAuthor_PutAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		id             int
		author         models.Author
		expectedOutput models.Author
		checkAuthor    bool
		err            error
	}{
		{
			"valid case", 1, models.Author{
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
			}, true, nil,
		},

		{
			"invalid id", 2, models.Author{
				ID:        2,
				FirstName: "gaurav",
				LastName:  "chandra",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, false, nil,
		},
		{
			"invalid", 100, models.Author{
				ID:        100,
				FirstName: "gaurav",
				LastName:  "chaudhari",
				Dob:       "18-07-2001",
				PenName:   "GCC",
			}, models.Author{}, true, errors.New("err"),
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockBookStore := datastore.NewMockBook(mockCtrl)
	mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
	ctx := context.Background()
	a := New(mockAuthorStore, mockBookStore)

	for _, v := range testcases {
		mockAuthorStore.EXPECT().CheckAuthorByID(ctx, v.id).Return(v.checkAuthor).AnyTimes()
		mockAuthorStore.EXPECT().PutAuthor(ctx, v.id, v.author).Return(v.expectedOutput, v.err).AnyTimes()

		output, _ := a.PutAuthor(ctx, v.id, v.author)
		if !reflect.DeepEqual(v.expectedOutput, output) {
			t.Errorf("Expected %v\tGot %v", v.expectedOutput, output)
		}
	}
}

func TestServiceAuthor_DeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		id          int
		deletedID   int
		checkAuthor bool
		errBook     error
		errAuthor   error
	}{
		{
			"valid case", 1, 1, true, nil, nil,
		},
		{
			"invalid case", 2, 0, true, errors.New("err"), nil,
		},
		{
			"valid case", 5, 0, true, nil, errors.New("err"),
		},
		{
			"valid case", 10, 0, false, nil, nil,
		},
	}
	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	for _, v := range testcases {
		mockBookStore := datastore.NewMockBook(mockCtrl)
		mockAuthorStore := datastore.NewMockAuthor(mockCtrl)
		ctx := context.Background()
		a := New(mockAuthorStore, mockBookStore)

		mockAuthorStore.EXPECT().CheckAuthorByID(ctx, v.id).Return(v.checkAuthor).AnyTimes()
		mockBookStore.EXPECT().DeleteBookByAuthorID(ctx, v.id).Return(v.errBook).AnyTimes()
		mockAuthorStore.EXPECT().DeleteAuthor(ctx, v.id).Return(v.deletedID, v.errAuthor).AnyTimes()

		id, _ := a.DeleteAuthor(ctx, v.id)
		if id != v.deletedID {
			t.Errorf("test case fail")
		}
	}
}

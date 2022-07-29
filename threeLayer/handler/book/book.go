package handlerbook

import (
	"strconv"
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"
)

type HandlerBook struct {
	serviceBook service.Book
}

func New(b service.Book) HandlerBook {
	return HandlerBook{
		serviceBook: b,
	}
}

// GetAllBooks  extract params from request and call GetAllBooks on  service layer to get all books
func (h HandlerBook) GetAllBooks(ctx *gofr.Context) (interface{}, error) {
	output, err := h.serviceBook.GetAllBooks(ctx)
	if err != nil {
		return output, err
	}

	return output, nil
}

// GetBookByID extract id and validate id  from request and call GetBookByID on service layer to get book detail
func (h HandlerBook) GetBookByID(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))

	if err != nil || id <= 0 {
		params := []string{ctx.Param("id")}
		return nil, errors.InvalidParam{Param: params}
	}

	var output models.Book

	output, err = h.serviceBook.GetBookByID(ctx, id)
	if err != nil {
		return nil, errors.EntityNotFound{}
	}

	return output, nil
}

// PostBook extract and validate book detail from request and call PostBook on service layer to post book
func (h HandlerBook) PostBook(ctx *gofr.Context) (interface{}, error) {
	var book models.Book
	if err := ctx.Bind(&book); err != nil {
		return nil, err
	}

	if !validateBook(&book) {
		return nil, errors.ForbiddenRequest{}
	}

	return h.serviceBook.PostBook(ctx, &book)
}

// PutBook extract and validate book detail from request and call PUtBook on service layer to update book
func (h HandlerBook) PutBook(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil || id <= 0 {
		return nil, errors.InvalidParam{}
	}

	var book models.Book

	err = ctx.Bind(&book)
	if err != nil {
		return nil, err
	}

	if !validateBook(&book) || !validateAuthor(book.Author) {
		return nil, errors.ForbiddenRequest{}
	}

	return h.serviceBook.PutBook(ctx, id, &book)
}

// DeleteBook extract and validate book id  from request and call DeleteBook on service layer to delete Book
func (h HandlerBook) DeleteBook(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil || id <= 0 {
		return nil, errors.InvalidParam{}
	}

	return h.serviceBook.DeleteBook(ctx, id)
}

// validateBook validate book detail
func validateBook(b *models.Book) bool {
	slc := strings.Split(b.PublicationDate, "-")
	sz := 3

	switch {
	case b.Publication != "Scholastic" && b.Publication != "Penguin" && b.Publication != "Arihanth":
		return false
	case len(slc) < sz:
		return false
	case slc[2] >= "2022" || slc[2] < "1880":
		return false
	case b.Title == "":
		return false
	default:
		return true
	}
}

// validateAuthor validate author
func validateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

func SetStatusCode(data interface{}, err error) error {
	switch err.(type) {
	case errors.EntityAlreadyExists:
		return errors.EntityAlreadyExists{}
	case errors.InvalidParam:
		return errors.InvalidParam{}
	case errors.EntityNotFound:
		return errors.EntityNotFound{}
	default:
		return nil
	}

}

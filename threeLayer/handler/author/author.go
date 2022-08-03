package handlerauthor

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"MyProject/CRUD-PROJECT/threeLayer/models"
	"MyProject/CRUD-PROJECT/threeLayer/service"
)

type HandlerAuthor struct {
	serviceAuthor service.Author
}

func New(author service.Author) HandlerAuthor {
	return HandlerAuthor{author}
}

// PostAuthor extract and validate author detail from request  and call PostAuthor service layer to post detail
func (h HandlerAuthor) PostAuthor(ctx *gofr.Context) (interface{}, error) {
	var author models.Author

	err := ctx.Bind(&author)
	if err != nil {
		return nil, err
	}

	if !ValidateAuthor(author) {
		return nil, errors.ForbiddenRequest{}
	}

	return h.serviceAuthor.PostAuthor(ctx, author)

}

// PutAuthor extract and validate author detail from request and call PutAuthor service layer to update detail
func (h HandlerAuthor) PutAuthor(ctx *gofr.Context) (interface{}, error) {
	var author models.Author

	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil || id <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = ctx.Bind(&author)
	if err != nil {
		return nil, err
	}

	if !ValidateAuthor(author) {
		return nil, errors.ForbiddenRequest{}
	}

	return h.serviceAuthor.PutAuthor(ctx, id, author)
}

// DeleteAuthor extract id and validate id and call DeleteAuthor on service layer
func (h HandlerAuthor) DeleteAuthor(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil || id <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return h.serviceAuthor.DeleteAuthor(ctx, id)
}

func ValidateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

package handlerauthor

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"
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

	id, err := strconv.Atoi(ctx.PathParam("id"))

	if err != nil || id <= 0 {
		params := []string{ctx.Param("id")}
		return nil, errors.InvalidParam{Param: params}
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
	id, err := strconv.Atoi(ctx.PathParam("id"))

	if err != nil || id <= 0 {
		params := []string{ctx.Param("id")}
		return nil, errors.InvalidParam{Param: params}
	}

	return h.serviceAuthor.DeleteAuthor(ctx, id)
}

func ValidateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

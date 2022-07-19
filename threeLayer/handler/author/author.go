package handlerauthor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"

	"github.com/gorilla/mux"
)

type HandlerAuthor struct {
	serviceAuthor service.Author
}

func New(author service.Author) HandlerAuthor {
	return HandlerAuthor{author}
}

// PostAuthor extract and validate author detail from request  and call PostAuthor service layer to post detail
func (h HandlerAuthor) PostAuthor(w http.ResponseWriter, req *http.Request) {
	var author models.Author

	err := UnMarshal(req, &author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateAuthor(author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	_, err = h.serviceAuthor.PostAuthor(ctx, author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// PutAuthor extract and validate author detail from request and call PutAuthor service layer to update detail
func (h HandlerAuthor) PutAuthor(w http.ResponseWriter, req *http.Request) {
	var author models.Author

	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = UnMarshal(req, &author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ValidateAuthor(author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	author, err = h.serviceAuthor.PutAuthor(ctx, id, author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = Marshal(w, author)
	if err != nil {
		return
	}
}

// DeleteAuthor extract id and validate id and call DeleteAuthor on service layer
func (h HandlerAuthor) DeleteAuthor(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	_, err = h.serviceAuthor.DeleteAuthor(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ValidateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

func Marshal(w http.ResponseWriter, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

func UnMarshal(req *http.Request, object interface{}) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &object)
	if err != nil {
		return err
	}

	return nil
}

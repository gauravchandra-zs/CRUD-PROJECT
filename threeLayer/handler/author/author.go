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

func (h HandlerAuthor) PostAuthor(w http.ResponseWriter, req *http.Request) {
	body, _ := io.ReadAll(req.Body)

	var author models.Author

	err := json.Unmarshal(body, &author)
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

func (h HandlerAuthor) PutAuthor(w http.ResponseWriter, req *http.Request) {
	var author models.Author

	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := io.ReadAll(req.Body)

	err = json.Unmarshal(body, &author)
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

	body, err = json.Marshal(author)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

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

package handlerauthor

import (
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

	_, err = h.serviceAuthor.PostAuthor(author)
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

	author, err = h.serviceAuthor.PutAuthor(id, author)
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

	_, _ = w.Write(body)
}

func (h HandlerAuthor) DeleteAuthor(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.serviceAuthor.DeleteAuthor(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

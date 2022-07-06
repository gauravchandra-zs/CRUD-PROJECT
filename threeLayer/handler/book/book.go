package handlerbook

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"Projects/GoLang-Interns-2022/threeLayer/models"
	"Projects/GoLang-Interns-2022/threeLayer/service"

	"github.com/gorilla/mux"
)

type HandlerBook struct {
	serviceBook service.Book
}

func New(b service.Book) HandlerBook {
	return HandlerBook{
		serviceBook: b,
	}
}

func (h HandlerBook) GetAllBooks(w http.ResponseWriter, req *http.Request) {
	params := map[string]string{
		"title":         req.URL.Query().Get("title"),
		"authorInclude": req.URL.Query().Get("includeAuthor"),
	}

	var output []models.Book

	output, err := h.serviceBook.GetAllBooks(params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		return
	}
}

func (h HandlerBook) GetBookById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	ID, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var output models.Book

	output, err = h.serviceBook.GetBookByID(ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (h HandlerBook) PostBook(w http.ResponseWriter, req *http.Request) {
	var book models.Book

	body, _ := io.ReadAll(req.Body)

	err := json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lastInsertedID, err := h.serviceBook.PostBook(&book)
	if err != nil || lastInsertedID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h HandlerBook) PutBook(w http.ResponseWriter, req *http.Request) {
	var book models.Book

	body, _ := io.ReadAll(req.Body)

	err := json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newData, err := h.serviceBook.PutBook(&book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err = json.Marshal(newData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h HandlerBook) DeleteBook(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deletedID, err := h.serviceBook.DeleteBook(id)
	if err != nil || deletedID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

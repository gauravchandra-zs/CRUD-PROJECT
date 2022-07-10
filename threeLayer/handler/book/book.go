package handlerbook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	ctx := context.Background()
	ctx = context.WithValue(ctx, "title", req.URL.Query().Get("title"))
	ctx = context.WithValue(ctx, "includeAuthor", req.URL.Query().Get("includeAuthor"))

	output, err := h.serviceBook.GetAllBooks(ctx)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h HandlerBook) GetBookById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var output models.Book

	ctx := context.Background()

	output, err = h.serviceBook.GetBookByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = marshal(w, output)
	if err != nil {
		return
	}
}

func (h HandlerBook) PostBook(w http.ResponseWriter, req *http.Request) {
	var book models.Book

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !validateBook(&book) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	_, err = h.serviceBook.PostBook(ctx, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h HandlerBook) PutBook(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !validateBook(&book) || !validateAuthor(book.Author) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	newData, err := h.serviceBook.PutBook(ctx, id, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = marshal(w, newData)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h HandlerBook) DeleteBook(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	_, err = h.serviceBook.DeleteBook(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

func validateAuthor(author models.Author) bool {
	if author.FirstName == "" || author.LastName == "" || author.Dob == "" || author.PenName == "" {
		return false
	}

	return true
}

func marshal(w http.ResponseWriter, book models.Book) error {
	body, err := json.Marshal(book)
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

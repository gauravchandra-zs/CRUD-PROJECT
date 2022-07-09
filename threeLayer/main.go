package main

import (
	"fmt"
	"log"
	"net/http"

	datastoreAuthor "Projects/GoLang-Interns-2022/threeLayer/datastore/author"
	datastoreBook "Projects/GoLang-Interns-2022/threeLayer/datastore/book"
	"Projects/GoLang-Interns-2022/threeLayer/drivers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	handlerAuthor "Projects/GoLang-Interns-2022/threeLayer/handler/author"
	handlerBook "Projects/GoLang-Interns-2022/threeLayer/handler/book"
	serviceAuthor "Projects/GoLang-Interns-2022/threeLayer/service/author"
	serviceBook "Projects/GoLang-Interns-2022/threeLayer/service/book"
)

func main() {
	db, err := drivers.CreateTable()
	if err != nil {
		log.Print(err)
		return
	}

	bookStore := datastoreBook.New(db)
	authorStore := datastoreAuthor.New(db)

	svcBook := serviceBook.New(bookStore, authorStore)
	svcAuthor := serviceAuthor.New(authorStore, bookStore)

	book := handlerBook.New(svcBook)
	author := handlerAuthor.New(svcAuthor)

	r := mux.NewRouter()
	r.HandleFunc("/book", book.GetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/book/{id}", book.GetBookById).Methods(http.MethodGet)
	r.HandleFunc("/book", book.PostBook).Methods(http.MethodPost)
	r.HandleFunc("/book/{id}", book.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book/{id}", book.PutBook).Methods(http.MethodPut)

	r.HandleFunc("/author", author.PostAuthor).Methods(http.MethodPost)
	r.HandleFunc("/author/{id}", author.DeleteAuthor).Methods(http.MethodDelete)
	r.HandleFunc("/author/{id}", author.PutAuthor).Methods(http.MethodPut)

	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	fmt.Println("server stared")

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

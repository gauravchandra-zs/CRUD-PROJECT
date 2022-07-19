package main

import (
	"log"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	datastoreAuthor "Projects/GoLang-Interns-2022/threeLayer/datastore/author"
	datastoreBook "Projects/GoLang-Interns-2022/threeLayer/datastore/book"
	"Projects/GoLang-Interns-2022/threeLayer/drivers"

	_ "github.com/go-sql-driver/mysql"

	handlerBook "Projects/GoLang-Interns-2022/threeLayer/handler/book"
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
	//	svcAuthor := serviceAuthor.New(authorStore, bookStore)

	book := handlerBook.New(svcBook)

	//r := mux.NewRouter()
	r := gofr.New()
	r.GET("/book", book.GetAllBooks)
	r.GET("/book/{id}", book.GetBookByID)
	r.POST("/book", book.PostBook)
	r.PUT("/book/{id}", book.PutBook)
	//r.HandleFunc("/book/{id}", book.DeleteBook).Methods(http.MethodDelete)
	//r.HandleFunc("/book/{id}", book.PutBook).Methods(http.MethodPut)
	//
	//r.HandleFunc("/author", author.PostAuthor).Methods(http.MethodPost)
	//r.HandleFunc("/author/{id}", author.DeleteAuthor).Methods(http.MethodDelete)
	//r.HandleFunc("/author/{id}", author.PutAuthor).Methods(http.MethodPut)

	//server := http.Server{
	//	Addr:    ":8000",
	//	Handler: r,
	//}
	//
	r.Start()
}

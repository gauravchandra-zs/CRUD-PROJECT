package main

import (
	"log"
	_ "net/http"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	datastoreAuthor "Projects/GoLang-Interns-2022/threeLayer/datastore/author"
	datastoreBook "Projects/GoLang-Interns-2022/threeLayer/datastore/book"
	"Projects/GoLang-Interns-2022/threeLayer/drivers"
	handlerauthor "Projects/GoLang-Interns-2022/threeLayer/handler/author"
	serviceauthor "Projects/GoLang-Interns-2022/threeLayer/service/author"

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
	svcAuthor := serviceauthor.New(authorStore, bookStore)

	book := handlerBook.New(svcBook)
	author := handlerauthor.New(svcAuthor)

	r := gofr.New()
	r.GET("/book", book.GetAllBooks)
	r.GET("/book/{id}", book.GetBookByID)
	r.POST("/book", book.PostBook)
	r.PUT("/book/{id}", book.PutBook)
	r.DELETE("/book/{id}", book.DeleteBook)

	r.POST("/author", author.PostAuthor)
	r.DELETE("/author/{id}", author.DeleteAuthor)
	r.PUT("/author/{id}", author.PutAuthor)

	r.Start()
}
